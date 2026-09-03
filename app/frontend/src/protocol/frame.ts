import type { MessageType } from "./opcode";

export interface Frame {
    type: MessageType;
    payload: Uint8Array;
}
