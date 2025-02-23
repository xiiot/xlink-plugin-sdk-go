package plugin

type Message struct {
	Kind     MessageKind       `yaml:"kind" json:"kind"`
	Metadata map[string]string `yaml:"meta" json:"meta"`
	Content  LazyValue         `yaml:"content" json:"content"`
}

type MessageKind string

const (
	// MessageReport device report message kind
	MessageDeviceReport MessageKind = "deviceReport"
	// MessageReport device lifecycle report message kind
	MessageDeviceState MessageKind = "deviceState"
)

type LazyValue struct {
	Value interface{}
	doc   []byte
}
