import asyncio
import nats
from nats.aio.client import Client, Subscription, Msg
from . import event
import pyee
import json
import logging

class Object:
    def __init__(self, nc: Client, id: str):
        self._id = id
        self._nc = nc
        self._sub: Subscription = None
        self.methods: dict = {}
        self.signals = pyee.EventEmitter()
        self.properties = pyee.EventEmitter()
        self.state = {}

    def id(self):
        return self._id

    async def subscribe(self):
        self._sub = await self._nc.subscribe(self._id, cb=self.msg_handler)

    async def unsubscribe(self):
        if self._sub:
            await self._sub.unsubscribe()

    ### Methods

    async def request_method(self, member: str, args: list) -> event.ObjectEvent:
        evt = event.new_invoke_event(self._id, member, args)
        reply = await self.request_event(evt)
        value = reply.value()
        return value
    
    def register_method(self, member: str, method):
        self.methods[member] = method
    

    ### Signals

    async def on_signal(self, member: str, cb: callable):
        self.signals.on(member, cb)

    async def publish_signal(self, member: str, args: list):
        evt = event.new_signal_event(self._id, member, args)
        await self.publish_event(evt)

    async def emit_signal(self, member: str, args: list):
        evt = event.new_signal_event(self._id, member, args)
        self.signals.emit(member, evt)

    ### Properties

    async def on_property(self, member: str, cb: callable):
        self.properties.on(member, cb)

    async def publish_property(self, member: str, value):
        evt = event.new_property_event(self._id, member, value)
        await self.publish_event(evt)

    async def emit_property(self, member: str, value):
        evt = event.new_property_event(self._id, member, value)
        self.properties.emit(member, evt)

    async def set_property(self, member: str, value):
        if self.state[member] != value:
            self.state[member] = value
            await self.emit_property(member, value)

    async def get_property(self, member: str):
        return self.state[member]
    
    async def get_properties(self):
        return self.state


    async def msg_handler(self, msg: Msg):
        logging.debug("Received a message on {subject}: {message}".format(
            subject=msg.subject, message=msg.data.decode()))
        evt = event.ObjectEvent.parse_raw(msg.data)
        if evt.event == event.EventTypeProperty:
            self.set_property(evt.member, evt.value())
        elif evt.event == event.EventTypeSignal:
            self.emit_signal(evt.member, evt.value())
        elif evt.event == event.EventTypeInvoke:
            if evt.member in self.methods:
                print("invoke", evt.member, evt.value())
                result = self.methods[evt.member](evt.value())
                reply = event.new_reply_event(self._id, evt.member, result)
                bytes = reply.json().encode()
                await msg.respond(bytes)
            else:
                print("No method", evt.member)
                reply = event.new_reply_event(self._id, evt.member, None)
                bytes = reply.json().encode()
                await msg.respond(bytes)



        print("Received a message on {subject}: {message}".format(
            subject=msg.subject, message=msg.data.decode()))
        
    async def publish_event(self, evt: event.ObjectEvent):
        return await self._nc.publish(self._id, evt.encode())
    
    async def request_event(self, evt: event.ObjectEvent) -> event.ObjectEvent:
        bytes = evt.json().encode()
        reply = await self._nc.request(self._id, bytes, timeout=1)
        return event.ObjectEvent.parse_raw(reply.data)
