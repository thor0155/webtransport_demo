import type { ChatState } from "@/models/chat-state";
import type { ChatHistoryResponse, ChatMessage } from "@/protocol/payload/chat";
import { ConnectionState } from "@/models/connect-state";
import type { WelcomePayload } from "@/protocol/payload/connection";
import type {
    JoinPayload,
    LeavePayload,
    MembersPayload,
    RoomInfoPayload,
} from "@/protocol/payload/room";
import type { ChatEventType } from "@/models/chat-event";
import type { LogLevel } from "@/models/log";

const MAX_EVENTS = 200;

export function applyLog(state: ChatState, level: LogLevel, message: string): ChatState {
    return {
        ...state,
        currentLog: { level, message },
    };
}

export function applyEventAndLog(state: ChatState, level: LogLevel, message: string): ChatState {
    return applyEvent(applyLog(state, level, message), level, message);
}

export function applyConnecting(
    state: ChatState,
    userId: string,
    userName: string,
    room: string,
): ChatState {
    return {
        ...state,

        connectionState: ConnectionState.Connecting,
        userId: userId,
        userName: userName,
        room,
        members: [],
        messages: [],
        events: [],
        chatHistory: {
            cursor: null,
            hasMore: true,
            loading: false,
        },
    };
}

export function applyReconnect(state: ChatState): ChatState {
    return {
        ...state,
        connectionState: ConnectionState.Reconnecting,
    };
}

export function applyEventAndReconnect(state: ChatState): ChatState {
    return applyEvent(applyReconnect(state), "reconnect", "reconnecting");
}

export function applyConnected(state: ChatState): ChatState {
    return {
        ...state,
        connectionState: ConnectionState.Connected,
    };
}

export function applyEventAndConnected(state: ChatState): ChatState {
    return applyEvent(applyConnected(state), "connect", "connected");
}
export function applyDisconnected(state: ChatState): ChatState {
    return {
        ...state,
        connectionState: ConnectionState.Disconnected,
    };
}

export function applyEventAndDisconnected(state: ChatState): ChatState {
    return applyEvent(applyDisconnected(state), "disconnect", "disconnected");
}

export function applyWelcome(state: ChatState, payload: WelcomePayload): ChatState {
    return {
        ...state,
        userId: payload.userId,
        session: payload.session,
    };
}

export function applyMembers(state: ChatState, payload: MembersPayload): ChatState {
    return {
        ...state,

        members: payload.members,
    };
}

export function applyJoin(state: ChatState, payload: JoinPayload): ChatState {
    const members = [...state.members];

    const index = members.findIndex((m) => m.id === payload.id);

    if (index >= 0) {
        members[index] = payload;
    } else {
        members.push(payload);
    }

    return {
        ...state,
        members,
    };
}

export function applyEventAndJoin(state: ChatState, payload: JoinPayload): ChatState {
    return applyEvent(applyJoin(state, payload), "join", `${payload.id} joined`);
}

export function applyLeave(state: ChatState, payload: LeavePayload): ChatState {
    return {
        ...state,

        members: state.members.map((member) =>
            member.id === payload.userId
                ? {
                      ...member,
                      online: false,
                  }
                : member,
        ),
    };
}

export function applyEventAndLeave(state: ChatState, payload: LeavePayload): ChatState {
    return applyEvent(applyLeave(state, payload), "leave", `${payload.userId} left`);
}

export function applyChat(state: ChatState, payload: ChatMessage): ChatState {
    return {
        ...state,
        messages: [
            ...state.messages,
            {
                id: payload.id,
                senderId: payload.senderId,
                senderName: payload.senderName,
                text: payload.text,
                timestamp: payload.timestamp,
            },
        ],
    };
}

export function applyEvent(state: ChatState, type: ChatEventType, message: string): ChatState {
    return {
        ...state,

        events: [
            ...state.events,
            {
                id: crypto.randomUUID(),
                type,
                message,
                timestamp: Date.now(),
            },
        ].slice(-MAX_EVENTS),
    };
}

export function applyClearEvents(state: ChatState): ChatState {
    return {
        ...state,
        events: [],
    };
}

export function applyRoomInfo(state: ChatState, payload: RoomInfoPayload): ChatState {
    return {
        ...state,
        roomInfo: payload,
    };
}

export function applyChatHistoryStart(state: ChatState): ChatState {
    return {
        ...state,
        chatHistory: {
            ...state.chatHistory,
            loading: true,
        },
    };
}

export function applyChatHistoryLoaded(state: ChatState, payload: ChatHistoryResponse): ChatState {
    const oldMessages = [...(payload.messages ?? [])].reverse();
    return {
        ...state,
        messages: [...oldMessages, ...state.messages],
        chatHistory: {
            cursor: payload.nextCursor,
            hasMore: payload.nextCursor !== null,
            loading: false,
        },
    };
}

export function applyChatHistoryError(state: ChatState): ChatState {
    return {
        ...state,
        chatHistory: {
            ...state.chatHistory,
            loading: false,
        },
    };
}
