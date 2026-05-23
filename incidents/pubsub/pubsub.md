# Asynchronous Messaging & Memory Outages

## The Roadmap Focus: Pub/Sub (PUBLISH, SUBSCRIBE) vs. Streams.

### The Incident: The Network Buffer Collapse

- Imagine Uber's dispatch system. When a user requests a ride, 
- the backend publishes the coordinates to a Redis Pub/Sub channel that all nearby drivers are subscribed to.

- **The Scenario:** 
  - Drivers drive into a tunnel. Their 4G/5G connections become incredibly slow or drop packets.

- **The Outage Mechanics:** 
  - Redis operates entirely in RAM. When it tries to push a Pub/Sub message to a driver's slow TCP socket, the OS network buffers fill up. 
  - This data backs up directly into Redis's memory.

- **The Limit Break:** 
  - Redis has a safeguard called client-output-buffer-limit pubsub. 
  - If a client's buffer exceeds this (usually 32MB), Redis violently disconnects the driver to protect itself.
  - When the driver leaves the tunnel and reconnects, the dispatch message is gone forever because Pub/Sub is fire-and-forget.

- **The OOM Death:** 
- If a junior DevOps engineer disables that limit to "stop drivers from disconnecting," 
- the slow network buffers will consume all available system RAM until the Linux OOM (Out of Memory) Killer terminates 
- the Redis process, bringing down the entire city's dispatch grid.

**The Solution: Trade-offs**
- **If using Pub/Sub:** 
  - You must accept data loss. 
  - Enforce strict client-output-buffer-limits, let drivers disconnect, and force the mobile app to do a REST API HTTP long-poll to fetch missed rides when they reconnect.

- If migrating to Streams:  
- Streams write to an append-only log in Redis. 
- If the driver goes offline, the data stays safe on the server. 
- When the driver reconnects, they use XREAD with their last known ID to fetch exactly what they missed. 

- The trade-off? Streams consume memory linearly and require aggressive trimming (XTRIM) to prevent disks from filling up.
- Go Implementation: Simulating a Pub/Sub Channel