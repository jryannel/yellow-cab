const nats = require('nats');
const { RemoteObject } = require('./object');

class Conn {
    constructor(addr) {
        this._addr = addr;
        this._objects = {};
        this._nc = null;
    }

    async connect() {
        this._nc = await nats.connect({ servers: this._addr });
    }

    async close() {
        this._nc.close();
    }

    newObject(id) {
        return new RemoteObject(this._nc, id);
    }

    object(id) {
        return this._objects[id];
    }

    ensureObject(id) {
        if (!this._objects[id]) {
            this._objects[id] = new RemoteObject(this._nc, id);
        }
        return this._objects[id];
    }

    async subscribe(id) {
        await this.ensureObject(id).subscribe();
    }

    async unsubscribe(id) {
        await this.ensureObject(id).unsubscribe();
    }

    async requestMethod(id, method, args) {
        return await this.ensureObject(id).requestMethod(method, args);
    }

    registerMethod(id, method, fn) {
        this.ensureObject(id).registerMethod(method, fn);
    }

    async publishSignal(id, signal, args) {
        return await this.ensureObject(id).publishSignal(signal, args);
    }


    async onSignal(id, signal, fn) {
        this.ensureObject(id).onSignal(signal, fn);
    }

    async offSignal(id, signal, fn) {
        this.ensureObject(id).offSignal(signal, fn);
    }

    async emitSignal(id, signal, args) {
        this.ensureObject(id).emitSignal(signal, args);
    }

    async onProperty(id, property, fn) {
        this.ensureObject(id).onProperty(property, fn);
    }

    async offProperty(id, property, fn) {
        this.ensureObject(id).offProperty(property, fn);
    }

    async emitProperty(id, property, args) {
        this.ensureObject(id).emitProperty(property, args);
    }

    setProperty(id, property, value) {
        this.ensureObject(id).setProperty(property, value);
    }

    getProperty(id, property) {
        return this.ensureObject(id).getProperty(property);
    }

    properties(id) {
        return this.ensureObject(id).properties();
    }
}

module.exports = {
    Conn,
};