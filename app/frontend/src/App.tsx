import { AppLayout } from "@/components/layout/AppLayout";
import { ChatPage } from "@/pages/ChatPage";
import { DashboardPage } from "@/pages/DashboardPage";
import { SettingsPage } from "@/pages/SettingsPage";
import type { AppPage } from "@/models/page";
import type { ReactNode } from "react";

function renderPage(page: AppPage): ReactNode {
    switch (page) {
        case "chat":
            return <ChatPage />;

        case "dashboard":
            return <DashboardPage />;

        case "settings":
            return <SettingsPage />;
    }
}

export default function App() {
    return <AppLayout>{(page) => renderPage(page)}</AppLayout>;
}
