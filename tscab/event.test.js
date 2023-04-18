const { ObjectEvent, Event, createInvokeEvent, createSignalEvent, createPropertyEvent } = require("./event")

test("ObjectEvent", () => {
    const e = new ObjectEvent("type", "id", "member", "data");
    expect(e.type).toBe("type");
    expect(e.id).toBe("id");
    expect(e.member).toBe("member");
    expect(e.data).toBe("data");
})

test("createPropertyEvent", () => {
    const e = createPropertyEvent("id", "member", "data");
    expect(e.type).toBe(Event.Property);
    expect(e.id).toBe("id");
    expect(e.member).toBe("member");
    expect(e.data).toBe("data");
})

test("createInvokeEvent", () => {
    const e = createInvokeEvent("id", "member", "data");
    expect(e.type).toBe(Event.Invoke);
    expect(e.id).toBe("id");
    expect(e.member).toBe("member");
    expect(e.data).toBe("data");
})

test("createSignalEvent", () => {
    const e = createSignalEvent("id", "member", "data");
    expect(e.type).toBe(Event.Signal);
    expect(e.id).toBe("id");
    expect(e.member).toBe("member");
    expect(e.data).toBe("data");
})

