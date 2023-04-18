package olink

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// PublishSignal
// EmitSignal
// OnSignal

// PublishProperty - send event
// EmitProperty - fire event
// OnProperty - register callback
// SetProperty (Set)
// PushProperty (Set + Publish)
// Property - get property value
// Properties - get all property values

type Object struct {
	lock       sync.RWMutex
	nc         *nats.Conn
	sub        *nats.Subscription
	id         string
	methods    map[string]MethodFunc
	properties *emitter[*ObjectEvent]
	signals    *emitter[*ObjectEvent]
	state      map[string]interface{}
}

func NewObject(nc *nats.Conn, id string) *Object {
	log.Info().Msgf("new object %s", id)
	return &Object{
		lock:       sync.RWMutex{},
		nc:         nc,
		id:         id,
		methods:    make(map[string]MethodFunc),
		properties: NewEmitter[*ObjectEvent](),
		signals:    NewEmitter[*ObjectEvent](),
		state:      make(map[string]interface{}),
	}
}

func (o *Object) Id() string {
	return o.id
}

// --------------------------------------------
// Subscription
// --------------------------------------------

func (o *Object) Subscribe() {
	if o.sub != nil {
		log.Info().Msgf("already subscribed to %s", o.id)
		return
	}
	sub, err := o.nc.Subscribe(o.id, func(msg *nats.Msg) {
		o.handleMessage(msg)
	})
	if err != nil {
		return
	}
	o.sub = sub
}

func (o *Object) Unsubscribe() {
	if o.sub != nil {
		o.sub.Unsubscribe()
		o.sub = nil
	}
}

// --------------------------------------------
// Methods
// --------------------------------------------

// RequestMethod sends a method request to the object and waits for the reply.
func (o *Object) RequestMethod(member string, args ...any) (any, error) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	evt := NewMethodEvent(o.id, member, args)
	reply, err := RequestEvent(o.nc, evt)
	if err != nil {
		return nil, err
	}
	value := reply.Value()
	return value, nil
}

// RegisterMethod registers a method handler.
func (o *Object) RegisterMethod(member string, fn MethodFunc) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.methods[member] = fn
}

// --------------------------------------------
// Signals
// --------------------------------------------

// OnSignal registers a signal handler.
func (o *Object) OnSignal(member string, fn func(evt *ObjectEvent)) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.signals.On(member, fn)
}

// PublishSignal publishes a signal event.
func (o *Object) PublishSignal(member string, args ...any) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	evt := NewSignalEvent(o.id, member, args)
	PublishEvent(o.nc, evt)
}

// EmitSignal emits a signal event.
func (o *Object) EmitSignal(member string, args ...any) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	evt := NewSignalEvent(o.id, member, args)
	o.signals.Emit(member, evt)
}

// --------------------------------------------
// Properties
// --------------------------------------------

// OnProperty registers a property handler.
func (o *Object) OnProperty(member string, fn func(evt *ObjectEvent)) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.properties.On(member, fn)
}

// PublishProperty publishes a property event.
func (o *Object) PublishProperty(member string, value interface{}) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	evt := NewPropertyEvent(o.id, member, value)
	PublishEvent(o.nc, evt)
}

// EmitProperty emits a property event.
func (o *Object) EmitProperty(member string, value interface{}) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	evt := NewPropertyEvent(o.id, member, value)
	o.properties.Emit(member, evt)
}

// SetProperty sets a property value and emits an event.
func (o *Object) SetProperty(member string, value interface{}) {
	o.lock.Lock()
	defer o.lock.Unlock()
	if o.state[member] == value {
		return
	}
	o.state[member] = value
	evt := NewPropertyEvent(o.id, member, value)
	o.properties.Emit(member, evt)
}

// Property gets a property value.
func (o *Object) Property(member string) interface{} {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.state[member]
}

// Properties gets all property values.
func (o *Object) Properties() map[string]interface{} {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.state
}

// --------------------------------------------
// Internal
// --------------------------------------------

// handleMessage handles a message.
func (o *Object) handleMessage(msg *nats.Msg) {
	var evt ObjectEvent
	err := json.Unmarshal(msg.Data, &evt)
	if err != nil {
		log.Error().Err(err).Msgf("failed to unmarshal object event: %s", msg.Subject)
		return
	}
	switch evt.Event {
	case EventTypeProperty:
		o.handleProperty(&evt)
	case EventTypeSignal:
		o.handleSignal(&evt)
	case EventTypeInvoke:
		o.handleInvoke(&evt, msg)
	default:
		log.Error().Msgf("invalid event: %v", evt)
	}
}

// handleProperty handles a property event.
func (o *Object) handleProperty(evt *ObjectEvent) {
	o.lock.Lock()
	defer o.lock.Unlock()
	v := evt.Value()
	o.SetProperty(evt.Member, v)
}

// handleSignal handles a signal event.
func (o *Object) handleSignal(evt *ObjectEvent) {
	o.lock.RLock()
	defer o.lock.RUnlock()
	args := evt.Args()
	o.EmitSignal(evt.Member, args)
}

// handleInvoke handles a method invoke request
func (o *Object) handleInvoke(evt *ObjectEvent, msg *nats.Msg) {
	if msg.Reply == "" {
		log.Error().Msgf("invalid method request: %s", msg.Subject)
		return
	}
	o.lock.RLock()
	fn, ok := o.methods[evt.Member]
	o.lock.RUnlock()
	if !ok {
		return
	}
	args := evt.Args()
	result, err := fn(args)
	if err != nil {
		log.Error().Err(err).Msgf("method failed: %s", evt.Member)
		return
	}
	evt.SetValue(result)
	bytes, err := json.Marshal(evt)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal method reply: %s", evt.Member)
		return
	}
	err = msg.Respond(bytes)
	if err != nil {
		log.Error().Err(err).Msgf("failed to send method reply: %s", evt.Member)
		return
	}
}

// --------------------------------------------
// Helpers
// --------------------------------------------

func PublishEvent(nc *nats.Conn, evt *ObjectEvent) {
	data, err := json.Marshal(evt)
	if err != nil {
		log.Error().Err(err).Msgf("failed to marshal event: %v", evt)
		return
	}
	err = nc.Publish(evt.ObjectId, data)
	if err != nil {
		log.Error().Err(err).Msgf("failed to publish event: %v", evt)
		return
	}
}

func RequestEvent(nc *nats.Conn, evt *ObjectEvent) (*ObjectEvent, error) {
	data, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}
	reply, err := nc.Request(evt.ObjectId, data, time.Second)
	if err != nil {
		return nil, err
	}
	var replyEvt ObjectEvent
	err = json.Unmarshal(reply.Data, &replyEvt)
	if err != nil {
		return nil, err
	}
	return &replyEvt, nil
}
