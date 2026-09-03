export interface ChatMessage {
    id: string;
    senderId: string;
    senderName: string;
    text: string;
    timestamp: number;
}

export interface ChatMessageRequest {
    room: string;
    text: string;
}

export interface ChatCursor {
    bucketSeq: number;
    msgId: string;
}

export interface ChatHistoryRequest {
    room: string;
    cursor: ChatCursor | null;
    limit: number;
}

export interface ChatHistoryResponse {
    messages: ChatMessage[];
    nextCursor: ChatCursor | null;
}
