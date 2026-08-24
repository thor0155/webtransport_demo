export class WelcomeService {
    private timer?: number;
    private readonly timeout: number;
    private readonly onTimeout: () => void;
    private completed: boolean = false;

    constructor(timeout: number, onTimeout: () => void) {
        this.timeout = timeout;
        this.onTimeout = onTimeout;
    }

    isCompleted(): boolean {
        return this.completed;
    }

    start() {
        this.stop();
        this.timer = window.setTimeout(() => {
            this.timer = undefined;
            this.onTimeout();
        }, this.timeout);
        this.completed = false;
    }

    complete() {
        this.stop();
        this.completed = true;
    }

    stop() {
        if (this.timer !== undefined) {
            clearTimeout(this.timer);
            this.timer = undefined;
        }
    }
}
