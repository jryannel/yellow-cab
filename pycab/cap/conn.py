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

    def new_object(self, objectId):
        o = object.Object(self.nc, objectId)
        self.objs[objectId] = o
    
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

    

