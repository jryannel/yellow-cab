import asyncio
import nats
from nats.aio.client import Client, Subscription
from . import object

class Conn:
    def __init__(self, addr: str):
        self.addr = addr
        self.nc: Client = None
        self.objs: dict[str, object.Object] = {}

    async def connect(self):
        self.nc = await nats.connect(self.addr)

    def connected_uri(self):
        return self.addr

    async def close(self):
        await self.nc.close()

    async def subscribe(self, objectId):
        return await self.ensure_object(objectId).subscribe()

    async def unsubscribe(self, objectId):
        return await self.ensure_object(objectId).unsubscribe()

    async def request_method(self, objectId, member, args):
        return await self.ensure_object(objectId).request_method(member, args)
    
    def register_method(self, objectId, member, method):
        return self.ensure_object(objectId).register_method(member, method)
    
    async def request_event(self, event):
        return await self.ensure_object(event.objectId).request_event(event)
    
    async def publish_signal(self, objectId, member, args):
        return await self.ensure_object(objectId).publish_signal(member, args)
    
    async def emit_signal(self, objectId, member, args):
        return await self.ensure_object(objectId).emit_signal(member, args)
    
    async def on_signal(self, objectId, member, cb):
        return await self.ensure_object(objectId).on_signal(member, cb)
    
    async def publish_property(self, objectId, member, value):
        return await self.ensure_object(objectId).publish_property(member, value)
    
    async def emit_property(self, objectId, member, value):
        return await self.ensure_object(objectId).emit_property(member, value)
    
    async def on_property(self, objectId, member, cb):
        return await self.ensure_object(objectId).on_property(member, cb)
    
    async def set_property(self, objectId, member, value):
        return await self.ensure_object(objectId).set_property(member, value)
    
    async def get_property(self, objectId, member):
        return await self.ensure_object(objectId).get_property(member)
    
    async def get_properties(self, objectId):
        return await self.ensure_object(objectId).get_properties()
    
    def get_object_ids(self):
        return self.objs.keys()
    
    def get_object(self, objectId):
        return self.objs[objectId]
    
    def register_object(self, objectId):
        return self.ensure_object(objectId)
    
    def unregister_object(self, objectId):
        if objectId in self.objs:
            del self.objs[objectId]
    
    def ensure_object(self, objectId):
        if objectId not in self.objs:
            self.objs[objectId] = object.Object(self.nc, objectId)            
        return self.objs[objectId]

    

