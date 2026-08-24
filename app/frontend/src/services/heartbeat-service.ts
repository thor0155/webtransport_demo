export interface HeartbeatOptions {
    interval: number;
    timeout: number;
}

const DefaultTimeout = 30000;
const DefaultInterval = 10000;

export class HeartbeatService {
    private pingTimer?: number;
    private timeoutTimer?: number;
    private readonly sendPing: () => void;
    private readonly options: HeartbeatOptions;
    private readonly onTimeout: () => void;

    constructor(options: HeartbeatOptions, sendPingFunc: () => void, onTimeoutFunc: () => void) {
        if (options.timeout <= 0) options.timeout = DefaultTimeout;
        if (options.interval <= 0) options.interval = DefaultInterval;

        this.options = options;
        this.sendPing = sendPingFunc;

        this.onTimeout = onTimeoutFunc;
    }

    start(): void {
        this.stop();
        this.schedulePing();
        this.scheduleTimeout();
    }

    stop(): void {
        if (this.pingTimer !== undefined) {
            clearInterval(this.pingTimer);
        }
        if (this.timeoutTimer !== undefined) {
            clearTimeout(this.timeoutTimer);
        }
    }

    pong(): void {
        this.scheduleTimeout();
    }

    private schedulePing(): void {
        this.pingTimer = window.setInterval(
            () => {
                this.sendPing();
            },

            this.options.interval,
        );
    }

    private scheduleTimeout(): void {
        if (this.timeoutTimer !== undefined) {
            clearTimeout(this.timeoutTimer);
        }

        this.timeoutTimer = window.setTimeout(() => {
            this.onTimeout();
        }, this.options.timeout);
    }
}
