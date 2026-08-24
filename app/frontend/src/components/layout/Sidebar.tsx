import type { AppPage } from "@/models/page";
import "./Sidebar.css";

interface Props {
    current: AppPage;

    onChange: (page: AppPage) => void;
}

const menus = [
    {
        id: "chat",
        label: "Chat",
        icon: "💬",
    },
    {
        id: "dashboard",
        label: "Dashboard",
        icon: "📊",
    },
    {
        id: "settings",
        label: "Settings",
        icon: "⚙️",
    },
] satisfies {
    id: AppPage;
    label: string;
    icon: string;
}[];

export function Sidebar({ current, onChange }: Props) {
    return (
        <aside className="sidebar">
            <div className="sidebar-header">WebTransport</div>

            <nav className="sidebar-menu">
                {menus.map((menu) => (
                    <button
                        key={menu.id}

                        className={current === menu.id ? "sidebar-item active" : "sidebar-item"}

                        onClick={() => onChange(menu.id)}
                    >
                        <span>{menu.icon}</span>

                        <span>{menu.label}</span>
                    </button>
                ))}
            </nav>
        </aside>
    );
}
