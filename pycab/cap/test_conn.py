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
    await c.subscribe("demo.calc")
    await c.unsubscribe("demo.calc")
    await c.close()

@pytest.mark.asyncio
async def test_request_method():
    addr = "nats://localhost:4222"
    c = conn.Conn(addr)
    await c.connect()
    await c.subscribe("demo.calc")
    def add(args):
        return sum(args)
    c.register_method("demo.calc", "add", add)
    result = await c.request_method("demo.calc", "add", [1, 2])
    assert result == 3
    await c.unsubscribe("demo.calc")
    await c.close()

