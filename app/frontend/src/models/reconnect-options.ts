export interface ReconnectOptions {
    enabled: boolean;

    maxAttempts: number;
    // unit ms
    initialDelay: number;
    // unit ms
    maxDelay: number;

    multiplier: number;

    onReconnect: () => Promise<void>;
    onReconnectStart: () => void;
}
