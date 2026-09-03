import type { LogLevel } from "@/models/log";
import type { ChatHistoryResponse, ChatMessage } from "@/protocol/payload/chat";
import type { WelcomePayload } from "@/protocol/payload/connection";
import type {
    MembersPayload,
    LeavePayload,
    JoinPayload,
    RoomInfoPayload,
} from "@/protocol/payload/room";

export interface ChatListener {
    onWelcome?(payload: WelcomePayload): void;

    onMembers?(payload: MembersPayload): void;

    onJoin?(payload: JoinPayload): void;

    onLeave?(payload: LeavePayload): void;

    onChat?(payload: ChatMessage): void;

    onChatHistory?(payload: ChatHistoryResponse): void;

    onPong?(): void;

    onConnected?(): void;

    onDisconnected?(error?: Error): void;

    onRoomInfo?(payload: RoomInfoPayload): void;

    onReconnect?(): void;

    onLog?(level: LogLevel, code: string, message: string): void;
}
