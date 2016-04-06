##Building a Chord Peer-to-Peer Overlay Network

###INTRODUCTION

A significant number of distributed hash table (DHT) models have been developed taking  peer-to-peer (P2P) applications explicitly into account. These prototypes include CAN, Chord, Pastry, Kademlia and more. But it is not true that DHT exclusively serves best P2P designs. By definition, P2P systems are merely networked environment where every node is responsible both as a client and a server. DHT-based infrastructure provides the overlay network with a certain way of locating data, since it can distribute content across nodes evenly, scale up the network size and route among the nodes fast. 

The purpose of this project is to build a DHT and utilize it in a P2P application. We chose Chord among the P2P-friendly DHT systems because it has a simple design and has been frequently studied and mentioned by researchers. Chord provides accessible material for implementing load balancing, scalability and routing efficiency. However, this early protocol does not yet elaborate the issues of correctness and reliability. P2P environments “have high churn rates: nodes join and leave the overlay continuously and do not stay in the overlay for long”. [Castro] In the context of DHT alone, there are possible cases that can go wrong during the network traffic. Studies have discovered and improved correctness-related invariants implicated in traditional DHT designs including Chord. Key challenge is at node failure or connection failure that the system is to maintain the topological structure.

This motivated problem is an on-going learning process and has been left as a soft requirement of this project. As a practical goal, this project aims to perform media streaming atop the final product. It requires querying the network for multiple files and serializing messages and data chunks. This shows that the ring-management properties of our Chord DHT implementation are valid, even though it has shortcomings. (Summarize more specific results about streaming.)


###DESIGN

This project carries the following functionalities. They are mostly based on the firstly presented Chord documentations [Stoica], unless specified otherwise. 

**Scalability:** Each node participates the network via links to its adjacent nodes. Theis local relation forms P2P overlays. When a node is added or removed, the adjacency state must be also coordinated locally. For example, a new node should get linked with a key; a removed node's key should be passed to a neighbor node. DHTs make this possible by storing the hash key for which this node is responsible. Changing node's responsibility to a certain key is done by reassigning the content of its hash table.

**Balanced Key Space:** Chord takes and extends the idea of consistent hashing that reduces total number of keys to reassign when a distributed system is resized. The reduction is achieved when each participating node only “views” its “roughly equal share”, not all of others. [Karger] As traditional hashing, node identifiers are randomized by a hash function.  Keys are then assigned to the “node whose identifier most immediately follows it” [Freedman]. While doing so, but in consistent hashing, the last node is glued to the first one. This gives illusion that the key space is now circular, moving in clockwise direction.  In this way, a node holds O(K/n) (where K = number of keys, n = number of nodes) keys. When this node leaves the system, the number of reassigned keys are only as many as the size of the individual hash table. 

**Routing Performance:** Above is still true in Chord, with table-size log(n).  The number of bits in a key is defined as the same value log(n). Using SHA1 hash function, a key is 160-bit long and a hash table becomes 160-keys long as well.  This means the DHT-based system now partitions the key  space so that each node knows 160 possible nodes to reach. Each hash table can be considered as a routing table storing ID-address pairs.  Note that it is also called finger table. The Chord DHT has every i-th key of the first node that succeeds or equals (n + 2i).  It is also possible to think that the i-th entry points the node in distance of (1 / 2n-i ). When querying a target node, a node recursively hops to the node of the largest key that is smaller than that of the target. Thus lookup efficiency or routing delay is O(log(n)). This greedy algorithm sets bound for required bandwidth which seems to work practically with large n. This improves poor performance of consistent hashing which possibly “requires computing all” the keys. [Karger] 

**Faulure Detection:**



From this design we assume the followings true:

* A node only has one successor and one predecessor.
* A node can reach any of other live nodes eventually.
* There is exactly one ring (no hole, no inside loops).


### RELATED WORK

There are studies that have interpreted above assumptions differently, or have articulated scenarios when they can get broken. 



### IMPLEMENTATION

The final product is developed using Go language.

1. For all nodes to agree on a single system state:

2. To utilize the Chord DHT in media streaming:


###REFERENCES

[Castro] M. Castro, M. Costa, and A. Rowstron. Performance and Dependability of Structured Peer-to-Peer Overlays. In Proc. of DSN, 2004.

[Stoica] I. Stoica, R. Morris, D. Karger, M. F. Kaashoek, and H. Balakrishnan. Chord: A Scalable Peer-to-peer Lookup Service for Internet Applications. In Proc. of SIGCOMM. ACM, 2001.

[Karger] D. Karger, E. Lehman, T. Leighton, Matthew Lewvine, D. Lewin, and R. Panigrahy.  Consistent Hashing and Random Trees: Distributed Caching Protocols for Relieving Hot Spots on the World Wide Web. In Proc. of the 29th Annual ACM Symposium on Theory of Computing, 1997.

[Freedman] M. J. Freedman, K. Lakshminarayana, S. Rhea, and I. Stoica. Non-Transitive Connectivity and DHTs. In Proc. of WORLDS, 2005
