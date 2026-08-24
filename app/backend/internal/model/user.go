package model

type User struct {
	Id   string `msgpack:"id"`
	Name string `msgpack:"name"`
}
