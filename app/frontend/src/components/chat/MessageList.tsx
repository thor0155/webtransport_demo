import { useEffect, useRef } from "react";

import { useChat } from "@/hooks/use-chat";

import { MessageItem } from "./MessageItem";

import "./MessageList.css";

export function MessageList() {
    const { state } = useChat();

    const bottomRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        bottomRef.current?.scrollIntoView({
            behavior: "smooth",
        });
    }, [state.messages.length]);

    return (
        <div className="message-list">
            {state.messages.length === 0 ? (
                <div className="message-empty">No messages yet</div>
            ) : (
                state.messages.map((message) => (
                    <MessageItem key={message.id} message={message} selfId={state.userId} />
                ))
            )}

            <div ref={bottomRef} />
        </div>
    );
}
