import {
    type MessageType,
    RequestMessageType,
    type RequestType,
    ResponseMessageType,
} from "@/protocol/opcode";

import { type Frame } from "@/protocol/frame";

import { encodePayload, decodePayload, encodeFrame } from "@/protocol/codec";

import type { HelloPayload, WelcomePayload } from "@/protocol/payload/connection";

import type {
    ChatCursor,
    ChatHistoryRequest,
    ChatHistoryResponse,
    ChatMessage,
    ChatMessageRequest,
} from "@/protocol/payload/chat";

import type {
    MembersPayload,
    JoinPayload,
    LeavePayload,
    RoomInfoPayload,
} from "@/protocol/payload/room";

import type { ConnectOptions } from "@/models/connect-options";
import type { Transport } from "@/transport/transport";
import { HeartbeatService } from "./heartbeat-service";
import type { ChatListener } from "./chat-listener";
import { ReconnectService } from "./reconnect-service";
import { WelcomeService } from "./welcome-service";
import type { TransportCloseEvent } from "@/models/transport-close-event";
import type { LogPayload } from "@/protocol/payload/log";

export class ChatService {
    private readonly transport: Transport;
    private readonly heartbeat: HeartbeatService;
    private readonly reconnect: ReconnectService;
    private readonly welcome: WelcomeService;
    private readonly listeners = new Set<ChatListener>();
    private lastConnectOptions?: ConnectOptions;
    private isConnected = false;

    private readonly handlers = new Map<MessageType, (frame: Frame) => void>([
        [ResponseMessageType.Welcome, this.handleWelcome.bind(this)],
        [ResponseMessageType.Members, this.handleMembers.bind(this)],
        [ResponseMessageType.Join, this.handleJoin.bind(this)],
        [ResponseMessageType.Leave, this.handleLeave.bind(this)],
        [ResponseMessageType.Chat, this.handleChat.bind(this)],
        [ResponseMessageType.ChatHistory, this.handleChatHistory.bind(this)],
        [ResponseMessageType.Pong, () => this.handlePong()],
        [ResponseMessageType.RoomInfo, this.handleRoomInfo.bind(this)],
        [ResponseMessageType.Log, this.handleLog.bind(this)],
    ]);

    constructor(transport: Transport) {
        this.transport = transport;
        this.transport.onMessage(this.handleFrame.bind(this));
        this.transport.onClose(this.handleTransportClosed.bind(this));
        this.heartbeat = new HeartbeatService(
            {
                interval: 30000,
                timeout: 60000,
            },
            this.sendPing.bind(this),
            this.handleHeartbeatTimeout.bind(this),
        );
        this.reconnect = new ReconnectService({
            enabled: true,
            maxAttempts: 5,
            initialDelay: 1000,
            maxDelay: 10000,
            multiplier: 2,
            onReconnect: () => this.reconnectInternal(),
            onReconnectStart: () => this.notify((listener) => listener.onReconnect?.()),
        });
        this.welcome = new WelcomeService(5000, this.handleWelcomeTimeout.bind(this));
    }

    addListener(listener: ChatListener) {
        this.listeners.add(listener);
    }

    removeListener(listener: ChatListener) {
        this.listeners.delete(listener);
    }

    async connect(option: ConnectOptions) {
        try {
            await this.transport.connect(option.transport);
            this.welcome.start();
            this.sendJoin(option.id, option.name, option.room);
            this.lastConnectOptions = option;
        } catch (err) {
            this.notifyError(err instanceof Error ? err.message : "Connect failed");
            throw err;
        }
    }

    sendGetHistory(room: string, cursor: ChatCursor | null, limit = 50) {
        const payload: ChatHistoryRequest = {
            room,
            cursor,
            limit,
        };
        this.send(RequestMessageType.ChatHistory, payload);
    }

    sendChat(roomName: string, message: string) {
        const payload: ChatMessageRequest = {
            text: message,
            room: roomName,
        };
        this.send(RequestMessageType.Chat, payload);
    }

    sendPing() {
        this.send(RequestMessageType.Ping);
    }

    sendJoin(id: string, name: string, room: string) {
        const payload: HelloPayload = {
            id,
            name,
            room,
        };
        this.send(RequestMessageType.Hello, payload);
    }

    disconnect() {
        this.reconnect.stop();
        this.closeTransport();
    }

    closeTransport() {
        this.leaveConnected();
        this.transport.close();
    }

    private handleWelcomeTimeout(): void {
        this.notifyError("Welcome timeout");

        this.closeTransport();

        this.reconnectLater();
    }

    private handleHeartbeatTimeout(): void {
        // console.warn("heartbeat timeout");
        this.notifyError("heartbeat timeout");
        this.closeTransport();
        this.reconnectLater();
    }

    private handleLog(frame: Frame) {
        const payload = decodePayload<LogPayload>(frame.payload);
        this.notify((listener) => listener.onLog?.(payload.level, payload.code, payload.message));
    }

    private handleWelcome(frame: Frame) {
        const payload = decodePayload<WelcomePayload>(frame.payload);

        this.enterConnected();
        this.notify((listener) => listener.onWelcome?.(payload));
    }

    private handleMembers(frame: Frame) {
        const payload = decodePayload<MembersPayload>(frame.payload);

        this.notify((listener) => listener.onMembers?.(payload));
    }

    private handleJoin(frame: Frame) {
        const payload = decodePayload<JoinPayload>(frame.payload);

        this.notify((listener) => listener.onJoin?.(payload));
    }

    private handleLeave(frame: Frame) {
        const payload = decodePayload<LeavePayload>(frame.payload);

        this.notify((listener) => listener.onLeave?.(payload));
    }

    private handleChat(frame: Frame) {
        const payload = decodePayload<ChatMessage>(frame.payload);

        this.notify((listener) => listener.onChat?.(payload));
    }

    private handlePong() {
        this.notify((listener) => listener.onPong?.());
        this.heartbeat.pong();
    }

    private handleFrame(frame: Frame) {
        const handler = this.handlers.get(frame.type);

        if (!handler) {
            console.warn("Unknown message:", frame.type);
            return;
        }

        handler(frame);
    }

    private handleRoomInfo(frame: Frame) {
        const payload = decodePayload<RoomInfoPayload>(frame.payload);

        this.notify((listener) => listener.onRoomInfo?.(payload));
    }

    private notify(callback: (listener: ChatListener) => void) {
        for (const listener of this.listeners) {
            try {
                callback(listener);
            } catch (error) {
                console.error(error);
            }
        }
    }

    private handleChatHistory(frame: Frame) {
        const payload = decodePayload<ChatHistoryResponse>(frame.payload);
        this.notify((listener) => listener.onChatHistory?.(payload));
    }

    private notifyError(message: string) {
        this.notify((listener) => listener.onLog?.("error", "", message));
    }

    private send<T>(type: RequestType, payload?: T): void {
        const bytes = encodeFrame({
            type,
            payload: payload === undefined ? new Uint8Array() : encodePayload(payload),
        });

        this.transport.write(bytes);
    }

    private async reconnectInternal() {
        if (!this.lastConnectOptions) {
            return;
        }

        try {
            await this.connect(this.lastConnectOptions);
        } catch (err) {
            // console.error("reconnect failed", err);
            this.notifyError(err instanceof Error ? err.message : "reconnect failed");
            this.reconnectLater();
        }
    }

    private handleTransportClosed(event: TransportCloseEvent) {
        this.disconnectInternal();

        switch (event.reason) {
            case "local":
                this.reconnect.stop();

                return;

            case "remote":
                this.reconnectLater();

                return;

            case "error":
                this.notifyError(event.error?.message ?? "Transport error");

                this.reconnectLater();

                return;
        }
    }

    private enterConnected() {
        if (this.isConnected) {
            return;
        }
        this.isConnected = true;
        this.welcome.complete();
        this.heartbeat.start();
        this.reconnect.reset();

        this.notify((listener) => listener.onConnected?.());
    }

    private leaveConnected() {
        if (!this.isConnected) {
            return;
        }
        this.isConnected = false;
        this.welcome.stop();
        this.heartbeat.stop();
    }

    private reconnectLater() {
        void this.reconnect.schedule();
    }

    private disconnectInternal() {
        this.leaveConnected();

        this.notify((listener) => listener.onDisconnected?.());
    }
}
