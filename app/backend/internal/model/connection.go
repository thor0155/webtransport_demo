package model

type WelcomeResponse struct {
	UserId  string `msgpack:"userId"`
	Session string `msgpack:"session"`
}

type HelloRequest struct {
	Id   string `msgpack:"id"`
	Name string `msgpack:"name"`
	Room string `msgpack:"room"`
}
