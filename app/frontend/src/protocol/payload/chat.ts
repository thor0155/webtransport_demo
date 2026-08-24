export interface ChatPayload {
    id: string;
    senderId: string;
    senderName: string;
    text: string;
    timestamp: number;
}

export interface ChatRequestPayload {
    text: string;
}
