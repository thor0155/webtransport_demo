import { useEffect, useMemo, useRef, useState } from "react";

import { ChatContext } from "./chat-context";

import { WebTransportClient } from "@/transport/webtransport";

import { ChatService } from "@/services/chat-service";

import type { ChatListener } from "@/services/chat-listener";

import type { ChatState } from "@/models/chat-state";

import type { WelcomePayload } from "@/protocol/payload/connection";

import type { MembersPayload, JoinPayload, LeavePayload } from "@/protocol/payload/room";

import type { ChatCursor, ChatHistoryResponse, ChatMessage } from "@/protocol/payload/chat";
import type { ConnectOptions } from "@/models/connect-options";
import { ConnectionState } from "@/models/connect-state";

import {
    applyChat,
    applyConnecting,
    applyDisconnected,
    applyEventAndConnected,
    applyEventAndDisconnected,
    applyEventAndLog,
    applyEventAndJoin,
    applyEventAndLeave,
    applyEventAndReconnect,
    applyMembers,
    applyRoomInfo,
    applyWelcome,
    applyClearEvents,
    applyChatHistoryStart,
    applyChatHistoryLoaded,
    applyChatHistoryError,
    // applyChatHistoryError,
} from "@/state/chat-state-actions";
import { CHAT_HISTORY_LIMIT } from "@/models/chat-message";
import { ErrorCode } from "@/protocol/error_opcode";

interface ChatProviderProps {
    children: React.ReactNode;
}

export function ChatProvider({ children }: ChatProviderProps) {
    const [state, setState] = useState<ChatState>({
        userId: "",
        connectionState: ConnectionState.Disconnected,
        room: "",
        members: [],
        messages: [],
        events: [],
        chatHistory: {
            cursor: null,
            hasMore: false,
            loading: false,
        },
    });

    const service = useMemo(() => new ChatService(new WebTransportClient()), []);
    const loadingChatHistoryLockRef = useRef(false);
    const welcomedRef = useRef(false);
    const currentRoomRef = useRef("");
    const chatCursorRef = useRef<ChatCursor | null>(null);

    const listener = useMemo<ChatListener>(
        () => ({
            onConnected() {
                setState(applyEventAndConnected);
            },

            onDisconnected() {
                loadingChatHistoryLockRef.current = false;
                welcomedRef.current = false;
                setState(applyEventAndDisconnected);
            },

            onWelcome(payload: WelcomePayload) {
                setState((prev) => applyWelcome(prev, payload));

                if (!welcomedRef.current) {
                    welcomedRef.current = true;
                    requestHistory(null);
                }
            },

            onMembers(payload: MembersPayload) {
                setState((prev) => applyMembers(prev, payload));
            },

            onJoin(payload: JoinPayload) {
                setState((prev) => applyEventAndJoin(prev, payload));
            },

            onLeave(payload: LeavePayload) {
                setState((prev) => applyEventAndLeave(prev, payload));
            },

            onChat(payload: ChatMessage) {
                setState((prev) => applyChat(prev, payload));
            },
            onChatHistory(payload: ChatHistoryResponse) {
                loadingChatHistoryLockRef.current = false;
                chatCursorRef.current = payload.nextCursor;
                setState((prev) => applyChatHistoryLoaded(prev, payload));
            },
            onRoomInfo(payload) {
                setState((prev) => applyRoomInfo(prev, payload));
            },

            onPong() {
                // heartbeat later
            },
            onReconnect() {
                setState(applyEventAndReconnect);
            },
            onLog(level, code, message) {
                if (code === ErrorCode.ChatHistoryFailure) {
                    loadingChatHistoryLockRef.current = false;
                    setState((prev) => applyChatHistoryError(prev));
                }
                setState((prev) => applyEventAndLog(prev, level, message));
            },
        }),
        [],
    );

    useEffect(() => {
        service.addListener(listener);
        return () => {
            service.removeListener(listener);
        };
    }, [service, listener]);

    async function connect(option: ConnectOptions) {
        welcomedRef.current = false;
        loadingChatHistoryLockRef.current = false;
        currentRoomRef.current = option.room;

        setState((prev) => applyConnecting(prev, option.id, option.name, option.room));

        try {
            await service.connect(option);
        } catch (error) {
            setState((prev) => applyDisconnected(prev));

            throw error;
        }
    }

    function disconnect() {
        service.disconnect();
        listener.onLeave?.({
            userId: state.userId,
            session: state.session!,
        });
        listener.onRoomInfo?.({
            room: state.room,
            memberCount: state.roomInfo ? state.roomInfo.memberCount - 1 : 0,
            maxMembers: state.roomInfo?.maxMembers ?? 0,
        });
    }

    function sendChat(roomName: string, message: string) {
        service.sendChat(roomName, message);
    }

    // 統一的實際送出函式,唯一真正呼叫 service.sendGetHistory 的地方
    function requestHistory(cursor: ChatCursor | null, limit = CHAT_HISTORY_LIMIT) {
        if (!welcomedRef.current) return; // Key point: I haven't received the "Welcome" message yet, so I just blocked it.
        if (loadingChatHistoryLockRef.current) return; // There are already requests in flight, please block duplicates.

        loadingChatHistoryLockRef.current = true;
        setState((prev) => applyChatHistoryStart(prev));
        service.sendGetHistory(currentRoomRef.current, cursor, limit);
    }

    // 給 MessageList 呼叫的對外函式,改成呼叫同一個內部函式
    function loadMoreHistory(limit: number = CHAT_HISTORY_LIMIT) {
        requestHistory(chatCursorRef.current, limit);
    }

    function clearEvents() {
        setState(applyClearEvents);
    }

    return (
        <ChatContext.Provider
            value={{
                state,
                connect,
                disconnect,
                sendChat,
                clearEvents,
                loadMoreHistory,
            }}
        >
            {children}
        </ChatContext.Provider>
    );
}
