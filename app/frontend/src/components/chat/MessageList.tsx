import { useEffect, useLayoutEffect, useRef } from "react";
import { useChat } from "@/hooks/use-chat";
import { MessageItem } from "./MessageItem";
import { DateDivider } from "./DateDivider";
import { groupMessagesByDay } from "@/utils/chattool";
import "./MessageList.css";

export function MessageList() {
    const { state, loadMoreHistory } = useChat();

    const containerRef = useRef<HTMLDivElement>(null);
    const topSentinelRef = useRef<HTMLDivElement>(null);
    const bottomRef = useRef<HTMLDivElement>(null);

    const prevScrollHeightRef = useRef(0);
    const prevMessageCountRef = useRef(0);
    const isAppendingHistoryRef = useRef(false);
    const hasScrolledToBottomRef = useRef(false);

    // 滑到頂端 sentinel 時觸發載入更多
    useEffect(() => {
        const container = containerRef.current;
        const sentinel = topSentinelRef.current;
        if (!container || !sentinel) return;

        const observer = new IntersectionObserver(
            (entries) => {
                if (
                    entries[0]?.isIntersecting &&
                    state.chatHistory.hasMore &&
                    !state.chatHistory.loading
                ) {
                    isAppendingHistoryRef.current = true;
                    prevScrollHeightRef.current = container.scrollHeight;
                    loadMoreHistory();
                }
            },
            { root: container, threshold: 0, rootMargin: "200px 0px 0px 0px" },
        );

        observer.observe(sentinel);
        return () => observer.disconnect();
    }, [state.chatHistory.hasMore, state.chatHistory.loading, loadMoreHistory]);

    // 訊息數量變化時的滾動處理
    useLayoutEffect(() => {
        const container = containerRef.current;
        if (!container) return;

        const countIncreased = state.messages.length > prevMessageCountRef.current;
        if (!countIncreased) {
            prevMessageCountRef.current = state.messages.length;
            return;
        }

        if (isAppendingHistoryRef.current) {
            // 插入舊訊息到頂端後,補償滾動位置避免畫面跳動
            const heightDiff = container.scrollHeight - prevScrollHeightRef.current;
            container.scrollTop += heightDiff;
            isAppendingHistoryRef.current = false;
        } else if (!hasScrolledToBottomRef.current) {
            bottomRef.current?.scrollIntoView({ behavior: "auto" });
            hasScrolledToBottomRef.current = true;
        } else {
            const distanceFromBottom =
                container.scrollHeight - container.scrollTop - container.clientHeight;
            if (distanceFromBottom < 150) {
                bottomRef.current?.scrollIntoView({ behavior: "smooth" });
            }
        }

        prevMessageCountRef.current = state.messages.length;
    }, [state.messages.length]);

    const groupedByDay = groupMessagesByDay(state.messages);

    return (
        <div className="message-list" ref={containerRef}>
            <div ref={topSentinelRef} className="top-sentinel" />

            {state.chatHistory.loading && <div className="message-status">loading...</div>}
            {!state.chatHistory.hasMore && state.messages.length > 0 && (
                <div className="message-status">The earliest message has been received.</div>
            )}

            {state.messages.length === 0 && !state.chatHistory.loading ? (
                <div className="message-empty">No messages yet</div>
            ) : (
                groupedByDay.map((group) => (
                    <div key={group.dateKey}>
                        <DateDivider label={group.label} />
                        {group.messages.map((message) => (
                            <MessageItem key={message.id} message={message} selfId={state.userId} />
                        ))}
                    </div>
                ))
            )}

            <div ref={bottomRef} />
        </div>
    );
}
