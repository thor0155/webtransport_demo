import type { ChatHistory, ChatMessage } from "./chat-message";
import type { Member } from "./member";
import type { ChatEvent } from "./chat-event";
import type { ConnectionState } from "./connect-state";
import type { RoomInfoPayload } from "@/protocol/payload/room";
import type { Log } from "./log";

export interface ChatState {
    connectionState: ConnectionState;
    userId: string;
    userName?: string;
    room: string;
    roomInfo?: RoomInfoPayload;
    session?: string;
    members: Member[];
    messages: ChatMessage[];
    events: ChatEvent[];
    currentLog?: Log;
    chatHistory: ChatHistory;
}
