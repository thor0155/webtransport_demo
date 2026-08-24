import { useChat } from "@/hooks/use-chat";
import { ConnectionState } from "@/models/connect-state";
import "./ChatHeader.css";

export function ChatHeader() {
    const { state } = useChat();

    function renderStatus() {
        switch (state.connectionState) {
            case ConnectionState.Connected:
                return {
                    text: "Connected",
                    className: "status connected",
                };

            case ConnectionState.Connecting:
                return {
                    text: "Connecting",
                    className: "status connecting",
                };

            case ConnectionState.Reconnecting:
                return {
                    text: "Reconnecting",
                    className: "status reconnecting",
                };

            case ConnectionState.Disconnecting:
                return {
                    text: "Disconnecting",
                    className: "status disconnecting",
                };

            default:
                return {
                    text: "Offline",
                    className: "status disconnected",
                };
        }
    }

    const status = renderStatus();

    const roomInfo = state.roomInfo;

    return (
        <header className="chat-header">
            <div className="header-info">
                <div className="header-item">
                    <span className="header-icon">💬</span>

                    <span>{state.room}</span>
                </div>

                <div className="header-item">
                    👥{" "}
                    <strong>
                        {roomInfo ? `${roomInfo.memberCount}/${roomInfo.maxMembers}` : "-"}
                    </strong>
                </div>

                <div className="header-item">{status.text}</div>
            </div>
        </header>
    );
}
