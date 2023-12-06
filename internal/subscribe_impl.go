package internal

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rizalarfiyan/be-revend/internal/handler"
)

type subscribe struct {
	conn mqtt.Client
}

func NewSubscribe(conn mqtt.Client) Subscribe {
	return &subscribe{
		conn: conn,
	}
}

func (s *subscribe) wait(token mqtt.Token) {
	token.Wait()
}

func (s *subscribe) BaseSubscribe(handler handler.MQTTHandler) {
	s.wait(s.conn.Subscribe("revend/trigger", 0, handler.Trigger))
}
