package plugin

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

type LazyValue struct {
	Value interface{}
	doc   []byte
}
