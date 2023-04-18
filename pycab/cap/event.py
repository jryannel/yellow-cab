import json

EventTypeProperty = "prop"
EventTypeSignal = "sig"
EventTypeInvoke = "inv"
EventTypeReply = "reply"
from pydantic import BaseModel

class ObjectEvent(BaseModel):
    event: str = None
    objectId: str = None
    member: str = None
    data: str = None

    def value(self):
        return json.loads(self.data)
    
    def set_value(self, value):
        self.data = json.dumps(value)

    def encode(self):
        return json.dumps(self.__dict__).encode()

def new_property_event(objectId: str, member: str, value) -> ObjectEvent:
    event = ObjectEvent()
    event.event = EventTypeProperty
    event.objectId = objectId
    event.member = member
    event.set_value(value)
    return event

def new_signal_event(objectId: str, member: str, value) -> ObjectEvent:
    event = ObjectEvent()
    event.event = EventTypeSignal
    event.objectId = objectId
    event.member = member
    event.set_value(value)
    return event

def new_invoke_event(objectId: str, member: str, value) -> ObjectEvent:
    event = ObjectEvent()
    event.event = EventTypeInvoke
    event.objectId = objectId
    event.member = member
    event.set_value(value)
    return event

def new_reply_event(objectId: str, member: str, value) -> ObjectEvent:
    event = ObjectEvent()
    event.event = EventTypeReply
    event.objectId = objectId
    event.member = member
    event.set_value(value)
    return event
