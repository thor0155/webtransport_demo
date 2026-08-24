import { useChat } from "@/hooks/use-chat";

import "./EventList.css";
import { useEffect, useRef } from "react";
import { EventCard } from "./EventCard";

export function EventList() {
    const { state, clearEvents } = useChat();
    const listRef = useRef<HTMLDivElement>(null);
    const bottomRef = useRef<HTMLDivElement>(null);
    const shouldStickRef = useRef(true);

    /**
     * 在底部
        ↓
        收到 Event
        ↓
        自動滑到底
     */
    useEffect(() => {
        if (!shouldStickRef.current) {
            return;
        }

        bottomRef.current?.scrollIntoView({
            behavior: "smooth",
            block: "end",
        });
    }, [state.events]);

    function handleScroll() {
        const element = listRef.current;

        if (!element) {
            return;
        }

        // 距離底部 16px 內，就視為使用者仍在看最新事件
        const distance = element.scrollHeight - element.scrollTop - element.clientHeight;

        shouldStickRef.current = distance < 16;
    }

    return (
        <section className="event-list">
            <header className="event-header">
                <h3 className="event-title">
                    System Events
                    <span className="event-count">({state.events.length})</span>
                </h3>

                <button
                    type="button"
                    className="event-clear"
                    onClick={clearEvents}
                    disabled={state.events.length === 0}
                >
                    Clear
                </button>
            </header>

            <div className="event-content" ref={listRef} onScroll={handleScroll}>
                {state.events.length === 0 ? (
                    <div className="event-empty">No events</div>
                ) : (
                    <>
                        {state.events.map((event) => (
                            <EventCard key={event.id} event={event} />
                        ))}

                        <div ref={bottomRef} />
                    </>
                )}
            </div>
        </section>
    );
}
