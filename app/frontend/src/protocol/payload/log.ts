import type { LogLevel } from "@/models/log";

export interface LogPayload {
    level: LogLevel;
    code: string;
    message: string;
}
