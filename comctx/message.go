package comctx

import (
	"encoding/json"
	"github.com/satori/go.uuid"
)

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

func NewRequestID() string {
	return uuid.NewV4().String()
}

type LazyValue struct {
	Value interface{}
	doc   []byte
}

func (v *LazyValue) Unmarshal(obj interface{}) error {
	if v.doc != nil {
		return json.Unmarshal(v.doc, obj)
	}
	if v.Value != nil {
		bs, err := json.Marshal(v.Value)
		if err != nil {
			return err
		}
		return json.Unmarshal(bs, obj)
	}
	return nil
}
