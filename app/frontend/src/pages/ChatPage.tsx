import "./ChatPage.css";

import { ChatHeader, EventList, MemberList, MessageInput, MessageList } from "@/components/chat";

import { Card } from "@/components/common";
import { ConnectionForm } from "@/components/connection/ConnectionForm";

export function ChatPage() {
    return (
        <div className="chat-page">
            <ChatHeader />

            <div className="chat-body">
                <section className="chat-main">
                    <Card title="Messages" className="message-card">
                        <MessageList />
                    </Card>

                    <MessageInput />
                </section>

                <aside className="chat-sidebar">
                    <Card title="Connection">
                        <ConnectionForm />
                    </Card>

                    <Card title="Members">
                        <MemberList />
                    </Card>

                    <Card title="Events">
                        <EventList />
                    </Card>
                </aside>
            </div>
        </div>
    );
}
