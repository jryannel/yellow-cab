const Event = {
    Property: 'prop',
    Invoke: 'inv',
    Signal: 'sig',
}

class ObjectEvent {
    constructor(type, id, member, data) {
        this.type = type;
        this.id = id;
        this.member = member;
        this.data = data;
    }
}

function createPropertyEvent(id, member, data) {
    return new ObjectEvent(Event.Property, id, member, data);
}

function createInvokeEvent(id, member, data) {
    return new ObjectEvent(Event.Invoke, id, member, data);
}

function createSignalEvent(id, member, data) {
    return new ObjectEvent(Event.Signal, id, member, data);
}

module.exports = {
    Event,
    ObjectEvent,
    createPropertyEvent,
    createInvokeEvent,
    createSignalEvent
};