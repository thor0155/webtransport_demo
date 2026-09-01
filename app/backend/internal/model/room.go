package model

type JoinPayload struct {
	Member
}

type LeavePayload struct {
	UserId  string `msgpack:"userId"`
	Session string `msgpack:"session"`
}

type RoomInfoPayload struct {
	Room        string `msgpack:"room"`
	MemberCount int    `msgpack:"memberCount"`
	MaxMembers  int    `msgpack:"maxMembers"`
}
