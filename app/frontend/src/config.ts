interface Config {
    apiBaseUrl: string;
    webtransportEndpoint: string;
    useCertHash: boolean;
    certHash: string;
}

export const config: Config = {
    apiBaseUrl: import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080",
    webtransportEndpoint: import.meta.env.VITE_WEBTRANSPORT_ENDPOINT ?? "https://localhost:8443",
    useCertHash: import.meta.env.VITE_WEBTRANSPORT_USE_CERT_HASH === "true",
    certHash: import.meta.env.VITE_WEBTRANSPORT_CERT_HASH,
};
