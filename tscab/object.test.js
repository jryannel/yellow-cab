const { Conn } = require('./conn');

test('RemoteObject', async () => {
    const c = new Conn('nats://localhost:4222');
    await c.connect();
    const o = c.newObject('demo.calc');
    await o.subscribe();
    console.log(o._id);
    expect(o.id()).toBe('demo.calc');
    expect(o.methods()).toEqual([]);
    expect(o.getProperties()).toEqual({});
    await c.close();
})

test('subscribe', async () => {
    const c = new Conn('nats://localhost:4222');
    await c.connect();
    const o = c.newObject('demo.calc');
    await o.subscribe();
    await o.unsubscribe();
    await c.close();
})


test('requestMethod', async () => {
    const c = new Conn('nats://localhost:4222');
    await c.connect();
    const o = c.newObject('demo.calc');
    await o.subscribe();
    o.registerMethod('add', (args) => {
        console.log('add', args);
        return args[0] + args[1];
    });
    const result = await o.requestMethod('add', [1, 2]);
    expect(result).toBe(3);
    await o.unsubscribe();
    await c.close();
})