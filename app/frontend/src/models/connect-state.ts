export const ConnectionState = {
    Disconnected: "disconnected",
    Connecting: "connecting",
    Connected: "connected",
    Disconnecting: "disconnecting",
    Reconnecting: "reconnecting",
} as const;

export type ConnectionState = (typeof ConnectionState)[keyof typeof ConnectionState];
