import type { ChatMessage } from "@/models/chat-message";
import { formatDateLabel } from "../date";

export interface MessageGroup {
    dateKey: string;
    label: string;
    messages: ChatMessage[];
}

export function groupMessagesByDay(messages: ChatMessage[]): MessageGroup[] {
    const groups: MessageGroup[] = [];
    let currentKey = "";

    for (const msg of messages) {
        const date = new Date(msg.timestamp);
        const dateKey = date.toDateString();

        if (dateKey !== currentKey) {
            currentKey = dateKey;
            groups.push({ dateKey, label: formatDateLabel(date), messages: [] });
        }
        groups[groups.length - 1]?.messages.push(msg);
    }
    return groups;
}
