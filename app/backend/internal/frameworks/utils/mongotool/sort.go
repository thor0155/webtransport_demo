package mongotool

import "go.mongodb.org/mongo-driver/v2/bson"

type SortDirection = int

const (
	ASC  SortDirection = 1
	DESC SortDirection = -1
)

type sortBuilder struct {
	sort bson.D
}

func NewSortBuilder() *sortBuilder {
	return &sortBuilder{}
}

func (s *sortBuilder) Add(key string, direction SortDirection) *sortBuilder {
	s.sort = append(s.sort, bson.E{Key: key, Value: direction})
	return s
}

func (s *sortBuilder) Build() bson.D {
	return s.sort
}
