
# The Concept: Redis Strings

In Redis, a String is the most basic data type. But don't let the name fool you—Redis strings are binary-safe,
meaning they can hold anything: plain text, raw JSON, integers, or even serialized binary data (like an image) up to 512 MB.

When you store an integer as a string, Redis is smart enough to understand it's a number, allowing you to use atomic commands like 
INCR or DECR without fetching, updating, and saving it back (which avoids race conditions).

To learn this hands-on, we will build a micro-feature: An API Rate Limiter combined with a Cache.

## Real-World Use Case: API Rate Limiter & Cache

- We will cache a user's API token details using SET and GET.

- We will enforce a dynamic expiration using EXPIRE and check the remaining time with TTL.

- We will track API hit counts using INCR.