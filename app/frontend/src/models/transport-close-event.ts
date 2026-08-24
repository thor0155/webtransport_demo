export interface TransportCloseEvent {
    reason: "local" | "remote" | "error";

    error?: Error;
}
