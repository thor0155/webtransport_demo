import "./MemberList.css";
import { useChat } from "@/hooks/use-chat";
import { MemberCard } from "./MemberCard";

export function MemberList() {
    const { state } = useChat();

    return (
        <div className="member-list">
            {state.members.length === 0 ? (
                <div className="member-empty">No members</div>
            ) : (
                state.members.map((member) => (
                    <MemberCard key={member.id} member={member} self={member.id === state.userId} />
                ))
            )}
        </div>
    );
}
