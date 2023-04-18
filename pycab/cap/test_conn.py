from . import conn
import nats
import pytest

@pytest.mark.asyncio
async def test_conn():
    addr = "nats://localhost:4222"
    c = conn.Conn(addr)
    await c.connect()
    assert c.connected_uri() == addr
    await c.close()

@pytest.mark.asyncio
async def test_subscribe():
    addr = "nats://localhost:4222"
    c = conn.Conn(addr)
    await c.connect()
    o = c.new_object("demo.calc")
    await o.subscribe("demo.calc")
    await o.unsubscribe("demo.calc")
    await c.close()

@pytest.mark.asyncio
async def test_request_method():
    addr = "nats://localhost:4222"
    c = conn.Conn(addr)
    await c.connect()
    o = c.new_object("demo.calc")
    await o.subscribe("demo.calc")
    def add(args):
        return sum(args)
    c.register_method("demo.calc", "add", add)
    result = await o.request_method("demo.calc", "add", [1, 2])
    assert result == 3
    await o.unsubscribe("demo.calc")
    await c.close()

