import type { ChatEventType } from "@/models/chat-event";

export function renderIcon(type: ChatEventType) {
    switch (type) {
        case "join":
            return "🟢";

        case "leave":
            return "🔴";

        case "connect":
            return "🟢";

        case "disconnect":
            return "⚪";

        case "reconnect":
            return "🔄";

        case "warn":
            return "⚠️";

        case "error":
            return "⛔";

        case "info":
            return "ℹ️";
    }
}
