import type { TransportOptions } from "@/models/connect-options";
import type { TransportCloseEvent } from "@/models/transport-close-event";
import type { Frame } from "@/protocol/frame";

export interface Transport {
    connect(options: TransportOptions): Promise<void>;

    close(): Promise<void> | void;

    write(frame: Uint8Array): void;

    onMessage(handler: (frame: Frame) => void): void;

    onClose(handler: (e: TransportCloseEvent) => void): void;
}
