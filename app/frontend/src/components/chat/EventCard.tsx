import type { ChatEvent } from "@/models/chat-event";
import { formatTime } from "@/utils/date";
import "./EventCard.css";
import { renderIcon } from "@/utils/event";

interface Props {
    event: ChatEvent;
}

export function EventCard({ event }: Props) {
    return (
        <div className={`event-card ${event.type}`}>
            <div className="event-icon">{renderIcon(event.type)}</div>

            <div className="event-body">
                <div className="event-message">{event.message}</div>

                <time className="event-time">{formatTime(event.timestamp)}</time>
            </div>
        </div>
    );
}
