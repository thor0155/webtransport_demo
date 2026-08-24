import { useContext } from "react";

import { ChatContext } from "@/providers/chat-context";

export function useChat() {
    const context = useContext(ChatContext);

    if (context === undefined) {
        throw new Error("useChat must be used inside ChatProvider");
    }

    return context;
}
