interface AvatarProps {
    name: string;
    online?: boolean;
}

export function Avatar({ name, online }: AvatarProps) {
    const text = name.slice(0, 2).toUpperCase();

    return (
        <div className="avatar">
            {text}
            {online && <span className="avatar-online" />}
        </div>
    );
}
