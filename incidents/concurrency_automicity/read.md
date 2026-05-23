# Concurrency & Atomicity (Transactions vs. Lua)

## The Roadmap Focus: WATCH, MULTI, EXEC, and Lua Scripting (EVAL).

### The Incident: The "Thundering Herd" Wallet Deduction

Imagine a Gojek Pay or Uber Cash scenario: A marketing push goes out, and 50,000 users try to claim a limited "50% OFF"
promo code (only 1,000 available) at the exact same millisecond.

If you read the remaining promo count, decrement it in Go, and save it back, you will hit massive race conditions.
Multiple Go routines will read Count: 5 and all write back Count: 4, giving away more codes than you actually have.

**Attempt 1: Optimistic Locking (WATCH / MULTI)** 
- You can use a Redis transaction. You WATCH the promo key. 
- If any other connection changes the key before you execute your transaction, your transaction fails.

The Scale Problem: 
- At 50,000 concurrent requests, 49,999 requests will fail the lock and have to retry. Your Go API
- servers will spike to 100% CPU purely from handling retry loops, causing a cascading failure (a Thundering Herd).

**Attempt 2: The Lua Script Solution (EVAL)** 

- Redis executes Lua scripts atomically. 
- Because Redis is single-threaded, while a Lua script is running, no other command
can run. 
- It guarantees perfect sequential consistency without the overhead of optimistic lock retries.

**The Production Trap:** 
- Because Lua scripts block the entire server, if you write a script with an O(N) time complexity 
- (e.g., iterating through a massive JSON array in Lua), Redis will freeze. 
- Health checks will fail, Kubernetes will assume the Pod is dead, and it will kill your Redis instance, causing a massive outage.
- Lua scripts must be incredibly fast and strictly O(1).

Go Implementation: Atomic Lua Wallet Transaction