package olink

import (
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type conn struct {
	lock    sync.RWMutex
	nc      *nats.Conn
	objects map[string]*Object
}

func NewConn(addr string) (*conn, error) {
	nc, err := nats.Connect(addr)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to NATS server")
		return nil, err
	}
	return &conn{
		lock:    sync.RWMutex{},
		nc:      nc,
		objects: make(map[string]*Object),
	}, nil
}

func (c *conn) Close() {
	for _, o := range c.objects {
		o.Unsubscribe()
	}
	c.nc.Close()
}

func (c *conn) ConnectedUrl() string {
	return c.nc.ConnectedUrl()
}

func (c *conn) NewObject(id string) *Object {
	c.lock.Lock()
	defer c.lock.Unlock()
	if o, ok := c.objects[id]; ok {
		return o
	}
	o := NewObject(c.nc, id)
	c.objects[id] = o
	return o
}

// --------------------------------------------
// Registry
// --------------------------------------------

// ObjectIds returns the object ids.
func (c *conn) ObjectIds() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	ids := make([]string, 0, len(c.objects))
	for id := range c.objects {
		ids = append(ids, id)
	}
	return ids
}

func (c *conn) Object(objectId string) *Object {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.objects[objectId]
}

func (c *conn) EnsureObject(objectId string) *Object {
	c.lock.RLock()
	defer c.lock.RUnlock()
	var o *Object
	if _, ok := c.objects[objectId]; ok {
		o = c.objects[objectId]
	} else {
		o = NewObject(c.nc, objectId)
		c.objects[objectId] = o
	}
	o.Subscribe()
	return o
}
