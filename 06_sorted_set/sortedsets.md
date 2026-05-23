# Sorted Sets/ZSets (ZADD, ZRANGE,ZINCRBY)

## The concept: Redis Sorted Sets 
 A Redis Sorted Set(ZSet) is one of the most powerful datastructures available in modern database
 
Like a reguler Set, every element is unique. 
However, every element is mapped to a floating-point number called *score*


Elements are maintained in a strictly sorted order based on their scores. 
If multiple elements share the same score, they are sorted lexicographically. 
This lets you retrieve ranges of elements by rank or score with incredible efficiency ($O(\log N)$ 
operation time via an underlying Skip List data structure).

## Real-World Use Case: Real-time Gaming Leaderboard 
 - We will build a live Gaming Leaderboard. Players will earn points dynamically,
 - and we will instantly fetch the top players alongside their ranks.