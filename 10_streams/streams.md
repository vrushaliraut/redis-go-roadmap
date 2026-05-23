# Streams (XADD, XREAD, XRANGE, XLEN)

## The Systemic Trade-Off & Scale Context

Before Redis 5.0, if you wanted append-only streaming logs with consumer tracking, 
you either pulled in heavy architecture like Apache Kafka or used Redis Lists (which lack multi-consumer fanout capabilities)
or Pub/Sub (which suffers from data loss if a consumer is offline).

Redis Streams function as a lightweight, persistent, append-only log file structure.

The Core Mechanism: It mirrors Apache Kafka architecture. 
It allows multiple consumers to listen to the same stream independently, allows data persistence (unlike Pub/Sub), 
and supports Consumer Groups so a pool of distributed Go microservices can load-balance message consumption 
without duplicate processing.

Hands-on Implementation: Transaction/Event Logging Stream