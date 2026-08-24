import { type Frame } from "@/protocol/frame";
// import { base64ToArrayBuffer } from "@/protocol/codec";
// import { config } from "@/config";
import { FrameReader } from "@/protocol/frame-reader";
import type { Transport } from "./transport";
import type { TransportOptions } from "@/models/connect-options";
import { withTimeout } from "@/utils/promise";
import type { TransportCloseEvent } from "@/models/transport-close-event";

export type MessageHandler = (frame: Frame) => void;

export class WebTransportClient implements Transport {
    private transport?: WebTransport;
    private writer?: WritableStreamDefaultWriter<Uint8Array>;
    private reader?: ReadableStreamDefaultReader<Uint8Array>;
    private messageHandler?: MessageHandler;
    private sendQueue: Uint8Array[] = [];
    private stopped = false;
    private localClosing = false;
    private closeHandler?: (e: TransportCloseEvent) => void;

    // -------------------------
    // Connect
    // -------------------------

    async connect(options: TransportOptions) {
        const wtOptions: WebTransportOptions = {};
        if (options.tls) {
            wtOptions.serverCertificateHashes = options.tls.serverCertificateHashes;
        }

        const transport = new WebTransport(options.url, wtOptions);

        try {
            await withTimeout(transport.ready, options.connectTimeout ?? 10000, "Connect timeout");
            this.transport = transport;
            void this.waitClosed();
        } catch (err) {
            this.closeTransport();
            throw err;
        }

        const stream = await transport.createBidirectionalStream();
        this.writer = stream.writable.getWriter();
        this.reader = stream.readable.getReader();
        this.stopped = false;
        this.localClosing = false;
        void this.readLoop();
        void this.writeLoop();
    }

    // -------------------------
    // Send API
    // -------------------------

    write(data: Uint8Array) {
        if (this.stopped) {
            throw new Error("Transport closed");
        }
        this.sendQueue.push(data);
    }

    // -------------------------
    // Register handler
    // -------------------------

    onMessage(handler: MessageHandler) {
        this.messageHandler = handler;
    }

    // -------------------------
    // Read Loop
    // -------------------------

    private async readLoop() {
        if (!this.reader) {
            return;
        }

        const frameReader = new FrameReader(this.reader);

        while (!this.stopped) {
            try {
                const frame = await frameReader.read();

                if (!frame) {
                    return;
                }

                this.messageHandler?.(frame);
            } catch {
                break;
            }
        }
    }

    // -------------------------
    // Write Loop
    // -------------------------

    private async writeLoop() {
        if (!this.writer) {
            return;
        }
        while (!this.stopped) {
            if (this.sendQueue.length === 0) {
                await this.sleep(5);
                continue;
            }

            const data = this.sendQueue.shift();

            if (!data) {
                continue;
            }

            try {
                await this.writer.write(data);
            } catch (err) {
                console.debug("writer closed", err);
                break;
            }
        }
    }

    // -------------------------
    // Close
    // -------------------------

    async close() {
        this.stopped = true;
        this.sendQueue.length = 0;

        try {
            await this.writer?.close();
        } catch (err) {
            console.debug("writer", err);
        } finally {
            this.writer = undefined;
        }

        try {
            await this.reader?.cancel();
        } catch (err) {
            console.debug("reader", err);
        } finally {
            this.reader = undefined;
        }

        try {
            this.closeTransport();
        } catch (err) {
            console.debug("transport", err);
        }
    }

    onClose(handler: (e: TransportCloseEvent) => void) {
        this.closeHandler = handler;
    }

    private sleep(ms: number) {
        return new Promise((resolve) => setTimeout(resolve, ms));
    }

    private closeTransport() {
        this.localClosing = true;
        this.transport?.close();
    }

    private async waitClosed() {
        const transport = this.transport;

        if (!transport) {
            return;
        }

        try {
            await transport.closed;

            if (this.localClosing) {
                this.closeHandler?.({
                    reason: "local",
                });
            } else {
                this.closeHandler?.({
                    reason: "remote",
                });
            }
        } catch (error) {
            this.closeHandler?.({
                reason: "error",
                error: error instanceof Error ? error : new Error("Transport closed"),
            });
        } finally {
            this.transport = undefined;
        }
    }
}
