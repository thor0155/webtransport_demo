export function formatTime(ts: number) {
    return new Date(ts).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
    });
}

export function formatDateLabel(date: Date): string {
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(today.getDate() - 1);
    if (date.toDateString() === today.toDateString()) return "today";
    if (date.toDateString() === yesterday.toDateString()) return "yesterday";
    return date.toLocaleDateString([], { year: "numeric", month: "long", day: "numeric" });
}
