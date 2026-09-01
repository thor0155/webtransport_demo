package model

import "time"

type Member struct {
	Id       string `msgpack:"id"`
	Name     string `msgpack:"name"`
	Session  string `msgpack:"session"`
	Online   bool   `msgpack:"online"`
	JoinedAt int64  `msgpack:"joinedAt"`
}

type Members struct {
	Members []*Member `msgpack:"members"`
}

func NewMember(sessionId string, user *User, online bool, joinedAt time.Time) Member {
	return Member{
		Id:       user.Id,
		Name:     user.Name,
		Session:  sessionId,
		Online:   online,
		JoinedAt: joinedAt.Unix(),
	}
}
