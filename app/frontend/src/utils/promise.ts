export function withTimeout<T>(
    promise: Promise<T>,
    timeout: number,
    message = "Operation timeout",
): Promise<T> {
    return Promise.race([
        promise,

        new Promise<T>((_, reject) => {
            const id = window.setTimeout(() => {
                reject(new Error(message));
            }, timeout);

            promise.finally(() => clearTimeout(id));
        }),
    ]);
}
