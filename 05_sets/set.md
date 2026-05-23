# Sets (SADD, SISMEMBER, SINTER)

## The concept: Redis Sets 
A Redis Set is an unordered collection of unique strings. If you try to add a duplicate item to srt, \
Redis simply ignores it. 

Because it behaves like a mathematical set, it operates in O(1) constant time for checking if an
item exists (SISMEMBER). Is also allows you to perform lightning-fast server-side operations like 
intersections(SINTER), unions (SUNION) and differences (SDIFF)

## Real-World Use Case: Unique Visitor Tracking & Social Graphs
- We will build an Analytics & Social Matching engine to:
- Track unique IP addresses visiting a web page (deduplication).
- Find common interests/friends between two users using set intersections.