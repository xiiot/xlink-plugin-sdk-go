package plugin

import "github.com/satori/go.uuid"

type Message struct {
	Kind     MessageKind       `yaml:"kind" json:"kind"`
	Metadata map[string]string `yaml:"meta" json:"meta"`
	Content  LazyValue         `yaml:"content" json:"content"`
}

type MessageKind string

const (
	MessageDeviceReport MessageKind = "deviceReport"
	MessageDeviceState  MessageKind = "deviceState"
)

const KeyDeviceName = "deviceName"

type LazyValue struct {
	Value interface{}
	doc   []byte
}

func NewRequestID() string {
	return uuid.NewV4().String()
}
