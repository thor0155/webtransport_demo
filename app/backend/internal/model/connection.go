package model

type WelcomePayload struct {
	UserId  string `msgpack:"userId"`
	Session string `msgpack:"session"`
}

type HelloPayload struct {
	Id   string `msgpack:"id"`
	Name string `msgpack:"name"`
	Room string `msgpack:"room"`
}
