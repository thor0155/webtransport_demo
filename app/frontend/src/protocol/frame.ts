import type { MessageType } from "./scheme";

export interface Frame {
    type: MessageType;
    payload: Uint8Array;
}
