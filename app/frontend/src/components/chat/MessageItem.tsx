import { formatTime } from "@/utils/date";
import "./MessageItem.css";

import type { ChatMessage } from "@/models/chat-message";

interface Props {
    message: ChatMessage;

    selfId: string;
}

export function MessageItem({ message, selfId }: Props) {
    const mine = message.senderId === selfId;

    return (
        <div className={mine ? "message-item mine" : "message-item"}>
            <div className="message-bubble">
                <div className="message-header">
                    <span className="message-sender">{message.senderName}</span>

                    <time>{formatTime(message.timestamp)}</time>
                </div>

                <div className="message-text">{message.text}</div>
            </div>
        </div>
    );
}
