export const ErrorCode = {
    ChatHistoryFailure: "GetHistoryFailure",
} as const;

export type ErrorCode = (typeof ErrorCode)[keyof typeof ErrorCode];
