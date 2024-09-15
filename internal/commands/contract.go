package commands

type Context interface {
	Send(what interface{}, opts ...interface{}) error
	ChatID() int64
}
