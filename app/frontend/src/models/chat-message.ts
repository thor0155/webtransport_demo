import type { ChatCursor } from "@/protocol/payload/chat";

export const CHAT_HISTORY_LIMIT: number = 50;

export interface ChatMessage {
    id: string;

    senderId: string;

    senderName: string;

    text: string;

    timestamp: number;
}

export interface ChatHistory {
    cursor: ChatCursor | null;
    hasMore: boolean;
    loading: boolean;
}
