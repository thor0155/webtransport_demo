import { useState } from "react";
import { useChat } from "@/hooks/use-chat";
import { config } from "@/config";
import type { ConnectOptions } from "@/models/connect-options";
import { ConnectionState } from "@/models/connect-state";
import { base64ToArrayBuffer } from "@/protocol/codec";
import { Button } from "@/components/common/Button";
import { GetCertHash } from "@/transport/http";

import "./ConnectionForm.css";
import { useMutation } from "@tanstack/react-query";

export function ConnectionForm() {
    const { state, connect, disconnect } = useChat();

    const [options, setOptions] = useState<ConnectOptions>({
        id: state.userId || "",
        name: state.userName || "",
        room: state.room || "",
        transport: {
            url: config.webtransportEndpoint,
            connectTimeout: 10000,
        },
    });

    const [verifyCertificate, setVerifyCertificate] = useState(config.useCertHash);

    const [certificateHash, setCertificateHash] = useState(config.certHash);

    const connected = state.connectionState === ConnectionState.Connected;

    const locked =
        state.connectionState === ConnectionState.Connecting ||
        state.connectionState === ConnectionState.Reconnecting ||
        state.connectionState === ConnectionState.Connected;

    const canConnect =
        options.id.trim() !== "" &&
        options.name.trim() !== "" &&
        options.room.trim() !== "" &&
        options.transport.url.trim() !== "";

    const getCertHashMutation = useMutation({
        mutationFn: GetCertHash,
        onSuccess: (data) => {
            setCertificateHash(data);
        },
        onError: (error) => {
            console.error("get cert hash failed:", error);
        },
    });

    function updateOption<K extends keyof ConnectOptions>(key: K, value: ConnectOptions[K]) {
        setOptions((prev) => ({
            ...prev,
            [key]: value,
        }));
    }

    function updateTransport(value: Partial<ConnectOptions["transport"]>) {
        setOptions((prev) => ({
            ...prev,

            transport: {
                ...prev.transport,
                ...value,
            },
        }));
    }

    function renderButtonState() {
        switch (state.connectionState) {
            case ConnectionState.Connecting:
                return "Connecting...";

            case ConnectionState.Reconnecting:
                return "Reconnecting...";

            case ConnectionState.Connected:
                return "Disconnect";

            default:
                return "Join";
        }
    }

    async function handleSubmit(event: React.SubmitEvent) {
        event.preventDefault();

        if (state.connectionState === ConnectionState.Reconnecting) {
            return;
        }

        if (connected) {
            disconnect();

            return;
        }

        await connect({
            ...options,

            transport: {
                ...options.transport,

                tls: verifyCertificate
                    ? {
                          serverCertificateHashes: [
                              {
                                  algorithm: "sha-256",
                                  value: base64ToArrayBuffer(certificateHash),
                              },
                          ],
                      }
                    : undefined,
            },
        });
    }

    return (
        <form className="connection-form" onSubmit={handleSubmit}>
            <fieldset disabled={locked}>
                <div className="form-section">
                    <h3>Identity</h3>

                    <div className="form-field">
                        <label>User ID</label>

                        <input
                            value={options.id}
                            onChange={(e) => updateOption("id", e.target.value)}
                        />
                    </div>

                    <div className="form-field">
                        <label>Name</label>

                        <input
                            value={options.name}
                            onChange={(e) => updateOption("name", e.target.value)}
                        />
                    </div>

                    <div className="form-field">
                        <label>Room</label>

                        <input
                            value={options.room}
                            onChange={(e) => updateOption("room", e.target.value)}
                        />
                    </div>
                </div>

                <div className="form-section">
                    <h3>Transport</h3>

                    <div className="form-field">
                        <label>URL</label>

                        <input
                            value={options.transport.url}

                            onChange={(e) =>
                                updateTransport({
                                    url: e.target.value,
                                })
                            }
                        />
                    </div>

                    <label className="checkbox-field">
                        <input
                            type="checkbox"

                            checked={verifyCertificate}

                            onChange={(e) => setVerifyCertificate(e.target.checked)}
                        />
                        Verify server certificate
                    </label>

                    <input
                        className="certificate-input"
                        disabled={!verifyCertificate}
                        value={certificateHash}
                        placeholder={"sha-256 base64 hash"}
                        onChange={(e) => setCertificateHash(e.target.value)}
                    />
                </div>
            </fieldset>

            <Button type="submit" fullWidth disabled={!canConnect}>
                {renderButtonState()}
            </Button>
            <Button
                type="button"
                onClick={() => getCertHashMutation.mutate()}
                disabled={getCertHashMutation.isPending}
            >
                {getCertHashMutation.isPending ? "getting..." : "Get Cert Hash"}
            </Button>
        </form>
    );
}
