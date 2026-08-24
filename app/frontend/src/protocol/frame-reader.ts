import { decodeFrame } from "./codec";
import { type Frame } from "./frame";
import { HEADER_SIZE } from "./header";

export class FrameReader {
    private reader: ReadableStreamDefaultReader<Uint8Array>;

    private buffer = new Uint8Array(0);

    constructor(reader: ReadableStreamDefaultReader<Uint8Array>) {
        this.reader = reader;
    }

    async read(): Promise<Frame | null> {
        while (true) {
            const frame = this.tryParse();

            if (frame) {
                return frame;
            }

            const result = await this.reader.read();

            if (result.done) {
                return null;
            }

            this.append(result.value);
        }
    }

    private tryParse(): Frame | null {
        if (this.buffer.length < HEADER_SIZE) {
            return null;
        }

        const view = new DataView(
            this.buffer.buffer,
            this.buffer.byteOffset,
            this.buffer.byteLength,
        );

        const length = view.getUint32(5, false);

        const totalSize = HEADER_SIZE + length;

        if (this.buffer.length < totalSize) {
            return null;
        }

        const frameBytes = this.buffer.slice(0, totalSize);

        this.buffer = this.buffer.slice(totalSize);

        return decodeFrame(frameBytes);
    }

    private append(chunk: Uint8Array) {
        const merged = new Uint8Array(this.buffer.length + chunk.length);

        merged.set(this.buffer, 0);

        merged.set(chunk, this.buffer.length);

        this.buffer = merged;
    }
}
