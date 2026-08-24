import { encode, decode } from "@msgpack/msgpack";
import { HEADER_SIZE, MAGIC, VERSION, type Header } from "./header";
import type { MessageType } from "./scheme";
import { FrameFlags } from "./flags";
import type { Frame } from "./frame";

export function encodePayload<T>(payload: T): Uint8Array {
    return encode(payload);
}

export function decodePayload<T>(data: Uint8Array): T {
    return decode(data) as T;
}

export function base64ToArrayBuffer(base64: string): ArrayBuffer {
    const binary = atob(base64);
    const bytes = new Uint8Array(binary.length);
    for (let i = 0; i < binary.length; i++) {
        bytes[i] = binary.charCodeAt(i);
    }
    return bytes.buffer;
}

function encodeHeader(header: Header): Uint8Array {
    const buffer = new Uint8Array(HEADER_SIZE);

    const view = new DataView(buffer.buffer);

    view.setUint16(0, header.magic, false);

    view.setUint8(2, header.version);

    view.setUint8(3, header.type);

    view.setUint8(4, header.flags);

    view.setUint32(5, header.length, false);

    return buffer;
}

function decodeHeader(buffer: Uint8Array): Header {
    if (buffer.length < HEADER_SIZE) {
        throw new Error("frame too small");
    }

    const view = new DataView(buffer.buffer, buffer.byteOffset, buffer.byteLength);

    return {
        magic: view.getUint16(0, false),

        version: view.getUint8(2),

        type: view.getUint8(3) as MessageType,

        flags: view.getUint8(4) as FrameFlags,

        length: view.getUint32(5, false),
    };
}

function validateHeader(header: Header): void {
    if (header.magic !== MAGIC) {
        throw new Error("invalid magic");
    }

    if (header.version !== VERSION) {
        throw new Error("unsupported version");
    }

    if (header.flags !== FrameFlags.None) {
        throw new Error("unsupported flags");
    }
}

export function encodeFrame(frame: Frame): Uint8Array {
    const header = encodeHeader({
        magic: MAGIC,

        version: VERSION,

        type: frame.type,

        flags: FrameFlags.None,

        length: frame.payload.length,
    });

    const result = new Uint8Array(HEADER_SIZE + frame.payload.length);

    result.set(header);

    result.set(frame.payload, HEADER_SIZE);

    return result;
}

export function decodeFrame(buffer: Uint8Array): Frame {
    const header = decodeHeader(buffer);

    validateHeader(header);

    const end = HEADER_SIZE + header.length;

    if (end > buffer.length) {
        throw new Error("invalid frame length");
    }

    return {
        type: header.type,

        payload: buffer.slice(HEADER_SIZE, end),
    };
}
