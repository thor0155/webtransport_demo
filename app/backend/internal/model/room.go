package model

type JoinResponse struct {
	Member
}

type LeaveResponse struct {
	UserId  string `msgpack:"userId"`
	Session string `msgpack:"session"`
}

type RoomInfoResponse struct {
	Room        string `msgpack:"room"`
	MemberCount int    `msgpack:"memberCount"`
	MaxMembers  int    `msgpack:"maxMembers"`
}
