const { EventEmitter } = require('eventemitter3');
const { createInvokeEvent, createSignalEvent, createPropertyEvent } = require('./event');
const nats = require('nats');

const sc = nats.StringCodec();

class RemoteObject {
    constructor(nc = null, id = null) {
        this._nc = nc;
        this._id = id;
        this._events = new EventEmitter();
        this._sub = null;
        this._methods = {};
        this._properties = new EventEmitter();
        this._signals = new EventEmitter();
        this._state = {};
    }

    id() {
        return this._id;
    }

    async subscribe() {
        this._sub = await this._nc.subscribe(this._id, {
            callback: (e, m) => {
                if (e) {
                    console.error(e);
                    return;
                }
                this._handleMessage(m);
            }
        });
    }

    async unsubscribe() {
        if (this._sub) {
            await this._sub.unsubscribe();
            this._sub = null;
        }
    }

    registerMethod(name, fn) {
        this._methods[name] = fn;
    }

    async requestMethod(method, args) {
        const event = createInvokeEvent(this._id, method, args);
        const reply = await this.requestEvent(event);
        return reply.data;
    }

    methods() {
        return Object.keys(this._methods);
    }

    async publishSignal(signal, args) {
        const event = createSignalEvent(this._id, signal, args);
        return await this.publishEvent(event);
    }

    async onSignal(signal, fn) {
        this._signals.on(signal, fn);
    }

    async offSignal(signal, fn) {
        this._signals.off(signal, fn);
    }

    async emitSignal(signal, args) {
        this._signals.emit(signal, args);
    }

    async onProperty(property, fn) {
        this._properties.on(property, fn);
    }

    async offProperty(property, fn) {
        this._properties.off(property, fn);
    }

    async emitProperty(property, value) {
        this._properties.emit(property, value);
    }

    setProperty(property, value) {
        if (this._state[property] === value) {
            return;
        }
        this._state[property] = value;
        this.emitProperty(property, value);
    }

    getProperty(property) {
        return this._state[property];
    }

    getProperties() {
        return this._state;
    }

    async publishProperty(property, value) {
        const event = createPropertyEvent(this._id, property, value);
        return await this.publishEvent(event);
    }


    async requestEvent(event) {
        const data = JSON.stringify(event);
        console.log('requestEvent', data);
        const msg = await this._nc.request(this._id, sc.encode(data), { timeout: 1000 });
        const reply = JSON.parse(sc.decode(msg.data));
        console.log('requestEvent response', reply);
        return reply
    }

    async publishEvent(event) {
        const data = JSON.stringify(event);
        await this._nc.publish(this._id, sc.encode(data));
    }

    async _handleMessage(msg) {
        const event = JSON.parse(sc.decode(msg.data));
        console.log('handleMessage', event);
        switch (event.type) {
            case 'prop':
                this.setProperty(event.member, event.data);
                break;
            case 'inv':
                const fn = this._methods[event.member];
                if (fn) {
                    const result = await fn(event.data);
                    const reply = createInvokeEvent(this._id, event.member, result);
                    const data = JSON.stringify(reply);
                    msg.respond(sc.encode(data));
                } else {
                    console.log('unknown method', event.member);
                }
                break;
            case 'sig':
                this.emitSignal(event.member, event.data);
                break;
            default:
                console.log('unknown event type', event.type);
        }
    }
}

module.exports = {
    RemoteObject,
};