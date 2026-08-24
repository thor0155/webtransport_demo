import type { FrameFlags } from "./flags";
import type { MessageType } from "./scheme";

export const MAGIC = 0x5754;
export const VERSION = 1;
export const HEADER_SIZE = 9;

export interface Header {
    magic: number;
    version: number;
    type: MessageType;
    flags: FrameFlags;
    length: number;
}
