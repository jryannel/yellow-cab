package olink

import (
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNewObject(t *testing.T) {
	c, err := NewConn(nats.DefaultURL)
	assert.Nil(t, err)
	defer c.Close()
	o := c.NewObject("test")
	assert.NotNil(t, o)
	assert.Equal(t, "test", o.Id())
}

func TestEnsureObject(t *testing.T) {
	c, err := NewConn(nats.DefaultURL)
	assert.Nil(t, err)
	defer c.Close()
	o := c.EnsureObject("test")
	assert.NotNil(t, o)
	assert.Equal(t, "test", o.Id())
}

func TestEnsureExisting(t *testing.T) {
	c, err := NewConn(nats.DefaultURL)
	assert.Nil(t, err)
	defer c.Close()
	o := c.NewObject("test")
	o2 := c.Object("test")
	assert.Equal(t, o, o2)
	o3 := c.EnsureObject("test")
	assert.Equal(t, o, o3)
}

func TestObjectIds(t *testing.T) {
	c, err := NewConn(nats.DefaultURL)
	assert.Nil(t, err)
	defer c.Close()
	o := c.NewObject("test")
	assert.Equal(t, []string{o.Id()}, c.ObjectIds())
}
