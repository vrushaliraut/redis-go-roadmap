# The Concept: Redis Lists

A Redis List is an ordered collection of strings sorted by insertion order. 
Under the hood, it is implemented as a Doubly Linked List.

Because it is a linked list, operations at the boundaries (adding or removing items from the very front or very back) are extremely fast—operating 
in $O(1)$ constant time—even if the list contains millions of elements. 
However, searching or accessing elements by index in the middle of the list takes $O(N)$ time.

## Real-World Use Case: Asynchronous Task Queue / Job Worker

Because elements can be pushed to one end and popped from the other, 
Redis Lists are universally used to build simple, ultra-fast Message Queues.
We will have a Producer push new background processing jobs (like sending emails) to the queue 
using LPUSH.

We will have a Worker poll and consume those jobs from the queue using RPOP 
(First-In, First-Out behavior).