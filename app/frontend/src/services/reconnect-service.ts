import type { ReconnectOptions } from "@/models/reconnect-options";

export class ReconnectService {
    private timer?: number;
    private attempts = 0;
    private readonly options: ReconnectOptions;

    constructor(options: ReconnectOptions) {
        this.options = options;
    }

    reset(): void {
        this.attempts = 0;

        this.stop();
    }

    stop(): void {
        if (this.timer !== undefined) {
            clearTimeout(this.timer);

            this.timer = undefined;
        }
    }

    async schedule(): Promise<void> {
        if (!this.options.enabled) {
            return;
        }

        if (this.attempts >= this.options.maxAttempts) {
            return;
        }

        const delay = this.getDelay();

        this.timer = window.setTimeout(
            async () => {
                this.attempts++;
                this.options.onReconnectStart();
                await this.options.onReconnect();
            },

            delay,
        );
    }

    private getDelay(): number {
        const delay = this.options.initialDelay * Math.pow(this.options.multiplier, this.attempts);

        return Math.min(delay, this.options.maxDelay);
    }
}
