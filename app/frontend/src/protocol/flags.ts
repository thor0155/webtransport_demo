export const FrameFlags = {
    None: 0,
} as const;

export type FrameFlags = (typeof FrameFlags)[keyof typeof FrameFlags];
