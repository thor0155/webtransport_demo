import type { Member } from "@/models/member";
import "./MemberCard.css";

interface Props {
    member: Member;
    self?: boolean;
}

export function MemberCard({ member, self = false }: Props) {
    return (
        <div className={self ? "member-card self" : "member-card"}>
            <div className="member-avatar">{member.name.slice(0, 1).toUpperCase()}</div>

            <div className="member-detail">
                <span className="member-name">{member.name}</span>

                {self && <span className="member-self">You</span>}
            </div>

            <span
                className={member.online ? "online-dot online" : "online-dot offline"}
                title={member.online ? "Online" : "Offline"}
            />
        </div>
    );
}
