export interface ChatPayload {
    id: string;
    senderId: string;
    senderName: string;
    text: string;
    timestamp: number;
}

export interface ChatRequestPayload {
    room: string;
    text: string;
}
