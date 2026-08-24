export interface ConnectOptions {
    id: string;
    name: string;
    room: string;
    transport: TransportOptions;
}

export interface TransportOptions {
    url: string;
    connectTimeout?: number;
    tls?: TLSOptions;
}

export interface TLSOptions {
    serverCertificateHashes: CertificateHash[];
}

export interface CertificateHash {
    algorithm: "sha-256";
    value: ArrayBuffer;
}
