package olink

import (
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestObjectRegister(t *testing.T) {
	c, err := NewConn(nats.DefaultURL)
	assert.NoError(t, err)
	defer c.Close()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	o := c.NewObject("demo.calc")
	defer o.Unsubscribe()
	o.Subscribe()
	o.RegisterMethod("add", func(args []any) (any, error) {
		return args[0].(float64) + args[1].(float64), nil
	})
	o.OnProperty("total", func(evt *ObjectEvent) {
		fmt.Printf("property: %s.%s = %v", evt.ObjectId, evt.Member, evt.Value())
	})
	o.OnSignal("clear", func(evt *ObjectEvent) {
		fmt.Printf("signal: %s.%s(%v)", evt.ObjectId, evt.Member, evt.Args())
	})
	result, err := o.RequestMethod("add", 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, 3.0, result)
}
