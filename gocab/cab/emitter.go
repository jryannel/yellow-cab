package olink

import "github.com/google/uuid"

func newId() string {
	return uuid.New().String()
}

type handler[T any] struct {
	id string
	fn func(evt T)
}

type emitter[T any] struct {
	handlers map[string][]handler[T]
}

func NewEmitter[T any]() *emitter[T] {
	return &emitter[T]{
		handlers: make(map[string][]handler[T]),
	}
}

func (e *emitter[T]) On(event string, fn func(evt T)) func() {
	id := newId()
	if _, ok := e.handlers[event]; !ok {
		e.handlers[event] = make([]handler[T], 0)
	}
	e.handlers[event] = append(e.handlers[event], handler[T]{id, fn})
	return func() {
		e.Off(event, id)
	}

}

func (e *emitter[T]) Off(event string, id string) {
	if _, ok := e.handlers[event]; !ok {
		return
	}
	for i, h := range e.handlers[event] {
		if h.id == id {
			e.handlers[event] = append(e.handlers[event][:i], e.handlers[event][i+1:]...)
			return
		}
	}
}

func (e *emitter[T]) Emit(event string, evt T) {
	if _, ok := e.handlers[event]; !ok {
		return
	}
	for _, h := range e.handlers[event] {
		h.fn(evt)
	}
}
