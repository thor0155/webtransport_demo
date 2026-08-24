import type { LogLevel } from "./log";

export type ChatEventType = "join" | "leave" | "connect" | "disconnect" | "reconnect" | LogLevel;

export interface ChatEvent {
    id: string;
    type: ChatEventType;
    message: string;
    timestamp: number;
}
