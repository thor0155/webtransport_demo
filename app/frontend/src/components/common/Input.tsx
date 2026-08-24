import "./Input.css";

import type { InputHTMLAttributes } from "react";

type InputProps = InputHTMLAttributes<HTMLInputElement>;

export function Input({ className, ...props }: InputProps) {
    return <input className={["input", className].filter(Boolean).join(" ")} {...props} />;
}
