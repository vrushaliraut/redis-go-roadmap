# Geospatial Indexes (GEOADD, GEOSEARCH)

## The Systemic Trade-Off & Scale Context
 
At Uber or Gojek scale, matching a passenger to the closest available drivers happens constantly. 
If you query a traditional database with relational latitude/longitude values using standard bounding-box calculations 
(SELECT WHERE lat BETWEEN...), index scanning falls short under high write loads.

Redis Geospatial Indexes store geographical coordinates natively by translating longitude
and latitude into a base-32 string using an algorithm called Geohash.

The Architecture Alignment: Under the hood,
Geospatial data is stored inside a Sorted Set (ZSet). 

The Geohash integer calculation serves as the score of the ZSet.

This turns a complex two-dimensional spatial proximity query into a highly optimized, single-dimensional $O(\log N + M)$ 
sorted range array scan.

Hands-on Implementation: Driver Proximity Search