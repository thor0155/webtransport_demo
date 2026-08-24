package model

import (
	"api/internal/utils"
	"time"
)

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

func EncodeJoinedMember(sessionId string, user *User, online bool, joinedAt time.Time) ([]byte, error) {

	m := JoinPayload{
		Member: NewMember(sessionId, user, online, joinedAt),
	}
	b, err := utils.Encode(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func EncodeLeavedMember(sessionId, userId string) ([]byte, error) {
	m := LeavePayload{
		UserId:  userId,
		Session: sessionId,
	}
	b, err := utils.Encode(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
