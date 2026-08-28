import { createContext } from "react";

import type { ChatState } from "@/models/chat-state";
import type { ConnectOptions } from "@/models/connect-options";

interface ChatContextValue {
    state: ChatState;

    connect(option: ConnectOptions): Promise<void>;

    disconnect(): void;

    sendChat(roomName:string, message: string): void;

    clearEvents(): void;
}

export const ChatContext = createContext<ChatContextValue | undefined>(undefined);
