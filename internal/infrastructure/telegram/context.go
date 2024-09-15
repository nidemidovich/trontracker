package telegram

import tele "gopkg.in/telebot.v3"

type Context struct {
	c tele.Context
}

func NewContext(c tele.Context) *Context {
	return &Context{c}
}

func (c *Context) Send(what interface{}, opts ...interface{}) error {
	return c.c.Send(what, opts...)
}

func (c *Context) ChatID() int64 {
	return c.c.Chat().ID
}
