import "./Card.css";

interface CardProps {
    title?: React.ReactNode;
    actions?: React.ReactNode;
    className?: string;
    children: React.ReactNode;
}

export function Card({ title, actions, className, children }: CardProps) {
    return (
        <section className={["card", className].filter(Boolean).join(" ")}>
            {(title || actions) && (
                <header className="card-header">
                    <div className="card-title">{title}</div>

                    <div className="card-actions">{actions}</div>
                </header>
            )}

            <div className="card-content">{children}</div>
        </section>
    );
}
