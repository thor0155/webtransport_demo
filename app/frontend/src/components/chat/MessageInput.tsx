import { useState, type KeyboardEvent } from "react";

import { useChat } from "@/hooks/use-chat";

import "./MessageInput.css";
import { ConnectionState } from "@/models/connect-state";
import { Button } from "../common/Button";

export function MessageInput() {
    const { state, sendChat } = useChat();
    const [message, setMessage] = useState("");
    const connected = state.connectionState === ConnectionState.Connected;

    function send() {
        const text = message.trim();

        if (!text || !connected) {
            return;
        }

        sendChat(text);

        setMessage("");
    }

    function handleKeyDown(event: KeyboardEvent<HTMLTextAreaElement>) {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();

            send();
        }
    }

    return (
        <div className="message-input">
            <textarea
                value={message}

                disabled={!connected}

                placeholder={connected ? "Type a message..." : "Connect first"}

                onChange={(event) => setMessage(event.target.value)}

                onKeyDown={handleKeyDown}
            />

            <Button disabled={!connected || !message.trim()} onClick={send}>
                Send
            </Button>
        </div>
    );
}
