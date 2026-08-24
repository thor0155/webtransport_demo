import "./Button.css";

import type { ButtonHTMLAttributes } from "react";

type ButtonVariant = "primary" | "secondary" | "danger";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: ButtonVariant;

    fullWidth?: boolean;
}

export function Button({
    variant = "primary",
    fullWidth = false,
    className = "",
    children,
    ...props
}: ButtonProps) {
    return (
        <button
            {...props}
            className={["button", `button-${variant}`, fullWidth && "button-full", className]
                .filter(Boolean)
                .join(" ")}
        >
            {children}
        </button>
    );
}
