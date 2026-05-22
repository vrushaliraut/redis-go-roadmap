# The Concept: Redis Hashes

A Redis *Hash* is a collection of field-value pairs, making it the perfect data structure to represent 
structured objects (like a User Profile, an E-commerce Product, or a Server Config).

While you could store a user profile as a flat JSON string inside a Redis String, 
Hashes are significantly more optimized:

1. Memory Efficiency: Redis compresses small hashes extremely well under the hood.

2. Partial Updates: You can update or fetch a single field (like changing a user's password) without having to read,
   deserialize, modify, and re-serialize the entire JSON object.

## Real-World Use Case: User Profile Management
We will build a simple User Profile Store. We'll save a user record with multiple fields, dynamically update a single field, check for existence, and pull the entire record. 