# Bitmaps (SETBIT, GETBIT, BITCOUNT, BITOP)

## The Systemic Trade-Off & Scale Context

- Imagine a high-scale application (like Gojek or Uber)
- tracking user daily active status or feature rollouts for 50 million monthly active users.
- Storing this as a list of user IDs or a boolean flag inside an explicit Hash table for every user consumes gigabytes
  of memory because of pointer overhead, hash collision chaining, and structural metadata.

- A Bitmap isn't an actual data structure; it's a string disguised as a bit array.
- Because Go strings and Redis strings are binary-safe, you can manipulate individual bits within a string.

- The Scale Match: 1 byte = 8 bits. If each bit represents a unique integer user_id offset (e.g., bit position 5,000,000
  corresponds to user ID 5,000,000), you can track the active status of 8 million users in just 1 Megabyte of memory.

- The Trap: If you have one user with ID 1 and another with ID 5,000,000, and nothing in between, Redis will still
  allocate all the bytes leading up to bit 5,000,000. It is highly efficient for dense integer IDs but bad for sparse
  keys (like UUIDs).

### Hands-on Implementation: High-Density Feature Flag Toggle / Daily Active Attendance