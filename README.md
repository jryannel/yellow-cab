# yellow-cab transport

An IPC framework based on NATS (https://nats.io) to create microservices using objects.

A object has an id and has properties, methods and signals. The properties are the state of the object. The methods are the actions that can be performed on the object. The signals are the events that can be emitted by the object.

Objects share a common NATS connection and are registered with the connection.


Typical object features:
- id
- subscribe / unsubscribe
- register methods
- remote request method and reply value
- register signal handler
- emit signal
- publish signal
- register property handler
- get property
- set property
- publish property
- emit property change

The connection typical features:
- connect to NATS server
- register object
- unregister object
- list registered objects
- get object


To run the microservices you need to run a nats-server. You can download it from https://nats.io/

## Code Generation (later)

The client shall be fully generated from the IDL. On the service side a adapter id generated to map the object events to typed function calls.

## Simulation / Monitor / Replay

The framework shall support a simulation mode where an simulation object can be registered. Monitoring should be done by simply subscribing to the NATS subject. The replay shall be possible using the NATS JetStream. Where each recoding is an own stream.


## Demos

We want to support the same demos for each language target to demonstrate the basic features of the framework.

### Counter
A counter object that can be incremented and decremented. The counter value is published as property.

Specific here is the propagation of properties and using several counter clients.

### Heater
A heater object which target temperature can be set. The heater object emits a signal when the target temperature is reached and also reports the current temperature and target temperature as properties. 

Specific here is the change of temperature over time and the signalling of the target temperature reached.


