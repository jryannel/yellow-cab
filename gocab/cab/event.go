package olink

import (
	"encoding/json"

	"github.com/rs/zerolog/log"
)

type MethodFunc func(args []any) (any, error)

type EventType string

const (
	EventTypeProperty EventType = "prop"
	EventTypeSignal   EventType = "sig"
	EventTypeInvoke   EventType = "inv"
	EventTypeReply    EventType = "reply"
)

type ObjectEvent struct {
	Event    EventType `json:"event"`
	ObjectId string    `json:"objectId"`
	Member   string    `json:"member"`
	Data     []byte    `json:"data"`
}

func (e *ObjectEvent) SetValue(value interface{}) {
	data, err := json.Marshal(value)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal object event: %s", e.Member)
		return
	}
	e.Data = data
}

func (e *ObjectEvent) Args() []interface{} {
	var args []interface{}
	err := json.Unmarshal(e.Data, &args)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal object event: %s", e.Member)
	}

	return args
}

func (e *ObjectEvent) KWArgs() map[string]interface{} {
	var kwargs map[string]interface{}
	err := json.Unmarshal(e.Data, &kwargs)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal object event: %s", e.Member)
	}
	return kwargs
}

func (e *ObjectEvent) Value() interface{} {
	var value interface{}
	err := json.Unmarshal(e.Data, &value)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal object event: %s", e.Member)
	}
	return value
}

func (e *ObjectEvent) Bytes(value interface{}) ([]byte, error) {
	e.SetValue(value)
	return json.Marshal(e)
}

func NewPropertyEvent(objectId string, member string, value interface{}) *ObjectEvent {
	evt := &ObjectEvent{
		Event:    EventTypeProperty,
		ObjectId: objectId,
		Member:   member,
	}
	evt.SetValue(value)
	return evt
}

func NewSignalEvent(objectId string, member string, args []interface{}) *ObjectEvent {
	evt := &ObjectEvent{
		Event:    EventTypeSignal,
		ObjectId: objectId,
		Member:   member,
	}
	evt.SetValue(args)
	return evt
}

func NewMethodEvent(objectId string, member string, args []interface{}) *ObjectEvent {
	evt := &ObjectEvent{
		Event:    EventTypeInvoke,
		ObjectId: objectId,
		Member:   member,
	}
	evt.SetValue(args)
	return evt
}

func NewReplyEvent(objectId string, member string, result interface{}) *ObjectEvent {
	evt := &ObjectEvent{
		Event:    EventTypeReply,
		ObjectId: objectId,
		Member:   member,
	}
	evt.SetValue(result)
	return evt
}

func DecodeEvent[T any](evt *ObjectEvent) (T, error) {
	var v T
	err := json.Unmarshal(evt.Data, &v)
	if err != nil {
		return v, err
	}
	return v, nil
}

func EncodeEvent[T any](evt *ObjectEvent, v *T) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	evt.Data = data
	return nil
}
