package codectool

import (
	"api/internal/model"
	"time"
)

func EncodeJoinedMember(user *model.User, sessionId string, online bool, joinedAt time.Time) ([]byte, error) {

	m := model.JoinResponse{
		Member: model.NewMember(sessionId, user, online, joinedAt),
	}
	b, err := Encode(m)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func EncodeLeavedMember(sessionId, userId string) ([]byte, error) {
	m := model.LeaveResponse{
		UserId:  userId,
		Session: sessionId,
	}
	b, err := Encode(m)
	if err != nil {
		return nil, err
	}
	return b, nil
}
