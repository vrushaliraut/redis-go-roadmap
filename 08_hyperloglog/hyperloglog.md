## HyperLogLog (PFADD, PFCOUNT, PFMERGE)

## - The Systemic Trade-Off & Scale Context 

Suppose you are designing a high-throughput endpoint for tracking unique search queries 
or unique ride views per city. 

If you use a Redis Set, your memory usage grows linearly ($O(N)$) with every single unique string added. 
At hundreds of millions of events, you will run out of RAM.

HyperLogLog (HLL) is a probabilistic data structure used to estimate the cardinality (unique count)of a
set.

The Magic: No matter how many millions or billions of unique items you add, every HyperLogLog key consumes a maximum
of 12 Kilobytes of memory.

The Trade-Off: It gives up absolute accuracy for massive space savings. It has a standard
error margin of 0.81%.

The Under the Hood Mechanics:
It hashes elements and inspects the distribution of leading zeros in
the binary format of those hashes to mathematically approximate how many unique elements would yield that
distribution.

Hands-on Implementation: Massive Scaled Unique View Counter