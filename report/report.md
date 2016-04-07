##Building a Chord Peer-to-Peer Overlay Network

###INTRODUCTION

A significant number of distributed hash table (DHT) models have been developed taking peer-to-peer (P2P) applications explicitly into account. These prototypes include CAN, Chord, Pastry, Kademlia and more. P2P is a general term for any networked environment where every node is responsible both as a client and a server. DHT provides a popular infrastructure for P2P as it locates data without a server. It also distributes content evenly, scales up the network size and routes among the nodes fast. Among those prototypes, we chose the Chord DHT because it has a simple design and has been frequently studied and mentioned by researchers. 

This project addressed the motivated problem of maintaining the Chord ring structure at node and communication failures. Regarding this concern, studies have discovered and improved possible incorrectness cases in traditional DHT designs including Chord.

Our interest is in approaching this known problem at the application level.  After building the DHT P2P overlay network, we have performed media streaming atop the system. This includes querying  serialized messages and data chunks. The results demonstrate two fold. First, it is a practical method to test the system. [ should elaborate – one sentence ]  Second, it may open possibility to detect and oversee the overlay topology from this higher layer, as a number of researchers suggested.


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


###RELATED WORK

There are studies that have interpreted above assumptions differently, or have articulated scenarios when they can get broken. 



----------------------------------------------------------------
Hey I am working on it right now & aim for writing UPTO this point. (including the abstract)
----------------------------------------------------------------


###IMPLEMENTATION

For this, I think it should be about how we integrate media streaming code into Chord DHT and BRBR... which I can't write about.


###EVALUATION

Same, it should be about the results from the streaming. 


###REFERENCES

[Castro] M. Castro, M. Costa, and A. Rowstron. Performance and Dependability of Structured Peer-to-Peer Overlays. In Proc. of DSN, 2004.

[Stoica] I. Stoica, R. Morris, D. Karger, M. F. Kaashoek, and H. Balakrishnan. Chord: A Scalable Peer-to-peer Lookup Service for Internet Applications. In Proc. of SIGCOMM. ACM, 2001.

[Karger] D. Karger, E. Lehman, T. Leighton, Matthew Lewvine, D. Lewin, and R. Panigrahy.  Consistent Hashing and Random Trees: Distributed Caching Protocols for Relieving Hot Spots on the World Wide Web. In Proc. of the 29th Annual ACM Symposium on Theory of Computing, 1997.

[Freedman] M. J. Freedman, K. Lakshminarayana, S. Rhea, and I. Stoica. Non-Transitive Connectivity and DHTs. In Proc. of WORLDS, 2005
