export class InvalidMagicError extends Error {
    constructor() {
        super("invalid frame magic");

        this.name = "InvalidMagicError";
    }
}

export class UnsupportedVersionError extends Error {
    constructor(version: number) {
        super(`unsupported protocol version: ${version}`);

        this.name = "UnsupportedVersionError";
    }
}

export class UnsupportedFlagsError extends Error {
    constructor(flags: number) {
        super(`unsupported frame flags: ${flags}`);

        this.name = "UnsupportedFlagsError";
    }
}

export class InvalidFrameLengthError extends Error {
    constructor() {
        super("invalid frame length");

        this.name = "InvalidFrameLengthError";
    }
}
