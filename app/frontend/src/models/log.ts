export type LogLevel = "info" | "warn" | "error";

export interface Log {
    level: LogLevel;
    message: string;
}
