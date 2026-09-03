export const RequestMessageType = {
    None: 0,
    Ping: 1,
    Hello: 2,
    Chat: 3,
    ChatHistory: 4,
} as const;

export const ResponseMessageType = {
    None: 0,
    Log: 1,
    Pong: 2,
    Leave: 3,
    Welcome: 4,
    Members: 5,
    Join: 6,
    Chat: 7,
    RoomInfo: 8,
    ChatHistory: 9,
} as const;

export type MessageType = number;
export type RequestType = (typeof RequestMessageType)[keyof typeof RequestMessageType];
