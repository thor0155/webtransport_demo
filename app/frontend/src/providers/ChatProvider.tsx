import { useEffect, useMemo, useState } from "react";

import { ChatContext } from "./chat-context";

import { WebTransportClient } from "@/transport/webtransport";

import { ChatService } from "@/services/chat-service";

import type { ChatListener } from "@/services/chat-listener";

import type { ChatState } from "@/models/chat-state";

import type { WelcomePayload } from "@/protocol/payload/connection";

import type { MembersPayload, JoinPayload, LeavePayload } from "@/protocol/payload/room";

import type { ChatPayload } from "@/protocol/payload/chat";
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
} from "@/state/chat-state-actions";

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
    });

    const service = useMemo(() => new ChatService(new WebTransportClient()), []);

    const listener = useMemo<ChatListener>(
        () => ({
            onConnected() {
                setState(applyEventAndConnected);
            },

            onDisconnected() {
                setState(applyEventAndDisconnected);
            },

            onWelcome(payload: WelcomePayload) {
                setState((prev) => applyWelcome(prev, payload));
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

            onChat(payload: ChatPayload) {
                setState((prev) => applyChat(prev, payload));
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
            onLog(level, message) {
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
        // console.log("connect options", option);

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

    function sendChat(message: string) {
        service.sendChat(message);
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
            }}
        >
            {children}
        </ChatContext.Provider>
    );
}
