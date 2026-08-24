import type { Member } from "@/models/member";

export type JoinPayload = Member;

export interface LeavePayload {
    userId: string;
    session: string;
}

export interface MembersPayload {
    members: Member[];
}

export interface RoomInfoPayload {
    room: string;
    memberCount: number;
    maxMembers: number;
    // owner?: string;
}
