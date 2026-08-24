import { useState } from "react";

import type { ReactNode } from "react";

import { Sidebar } from "./Sidebar";

import type { AppPage } from "@/models/page";
import "./AppLayout.css";

interface Props {
    children: (page: AppPage) => ReactNode;
}

export function AppLayout({ children }: Props) {
    const [page, setPage] = useState<AppPage>("chat");

    return (
        <div className="app-layout">
            <Sidebar
                current={page}

                onChange={setPage}
            />

            <main className="main-content">{children(page)}</main>
        </div>
    );
}
