# QueueToNull

QueueToNull is a REST API that accepts requests of different kinds, which is posted to RabbitMQ, then read by background processes from each queue. The project is purely made in order to explore possible design patterns for RabbitMQ, and as such does very little actual processing. Most of the examples are lifted directly from the [RabbitMQ Go tutorial](https://www.rabbitmq.com/tutorials/tutorial-one-go), though implemented in the context of a REST API application.

## Patterns

- [x] "Hello, World!"
- [x] Simple Work Queues
- [x] Publish/Subscribe
- [ ] Routing
- [ ] Topics
- [ ] RPC
