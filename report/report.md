#Streaming on Chord

###INTRODUCTION

A significant number of distributed hash table (DHT) models have been designed with intent to be used in peer-to-peer (P2P) applications. These prototypes include CAN, Chord, Pastry, Kademlia and more. P2P is a general term for any networked environment where every node is responsible both as a client and a server. DHT provides a popular infrastructure for P2P as it locates data without a server. It also scales the network up, distributes content evenly and routes among the nodes fast. Among those prototypes, we chose the Chord DHT because it has a simple design and has been frequently studied and mentioned by researchers. 

This project has addressed the motivated problem of maintaining the Chord ring overlay structure at node and communication failures. Regarding this concern, studies have discovered and improved possible incorrectness cases in traditional DHT designs including Chord. 

Our interest is in approaching this known problem at the application level.  After building a DHT P2P system, we have performed media streaming atop the system. This involves querying  serialized messages and data chunks. The results demonstrate two fold. First, it is a practical method to test the system. [ should elaborate – one sentence ] Second, it may open possibility to detect and oversee the overlay topology, as a number of researchers suggested.


###DESIGN

The following functionalities are based on settings and assumptions documented in the firstly presented Chord paper [Stoica], unless specified otherwise.

**Scalability:** Each node participates the network via links to its adjacent nodes. This relation forms P2P overlays. When a node is added or removed, the adjacency state is coordinated locally. That is, a node only knows the keys to nodes that this node is responsible for. In this context, a removed node's keys should be passed to a neighbor node; a new node should own some keys. DHTs make these possible by reassigning local hash tables in the way distributing the responsibilities.

**Balanced Key Space:** Chord takes and extends the idea of consistent hashing that reduces number of keys to reassign when a distributed system is resized. The reduction is achieved when each participating node only “views” its “roughly equal share”. [Karger]  As traditional hashing, nodes are identified by keys randomized by a hash function.  The keys are then assigned to the node whose key “most immediately follows it” [Freedman]. While doing so, but in consistent hashing, the last node is glued to the first one. This gives illusion that the key space is now circular, moving in clockwise direction.  In this way a node on average holds O(K/n) keys, where K = number of keys, n = number of nodes. When this node leaves the system, the number of reassigned keys is only as large as the size of the local hash table. 

**Routing Performance:** Above is still true in Chord, with table-size log n.  The number of bits in a key is defined as the same value log(n). Using SHA1 hash function, a key is 160-bit long and a hash table becomes 160-keys long as well.  This means the DHT-based system now partitions the key space so that each node knows 160 possible nodes to reach. Each hash table can be considered as a routing table storing key-address pairs (i.e. fingers).  The Chord DHT has every i-th key of the first node that succeeds or equals (n + 2i).  The i-th entry then points the node in distance of (1 / 2n-i). When querying a target node, a node recursively hops to the node of the largest key that is smaller than that of the target. This ensures routing delay within O(log n). This greedy algorithm sets bound for required bandwidth which seems to work practically with large n. This is an enhancement from consistent hashing which possibly “requires computing all” the keys. [Karger] 

**Node subscription and exit/failure:** To add or remove a node, two tasks should get done: updating all the routing tables that involve the node; and, transferring key responsibility from/to its new successor.  These operations are known to cost time complexity O(log n)^2.  Assume node B is node A's successor. To add C,  the system updates the predecessor A's routing table according to C, and recursively do the same through the routing tables of all other nodes that have entries of B. C gains the responsibility for keys in the range of (a, c] which was previously part of B's key (a, b]. Now B is only responsible for keys in (c, b]. (Figure 1) Similarly, when C leaves or gets detected to be failed,  all keys stored in C's routing table become reassigned to it's successor B.  All  associated routing table removes C. The successor B now owns C's keys as well as its previous keys. (Figure 2)

<img border=0 src='figure1.png'>

> Figure 1 (left).   Key transfer with new node.

> Figure 2 (right).    Key transfer with removed node.

It is important to keep track of nodes' state so that the Chord system keeps invariants true for data consistency. That is, “all members agree about which members store values for which keys.” [Zave] This topic is beyond the scope of this project. Rather, we narrow our focus down into how a node can reach its adjacent overlay peers correctly. A node's reachability is a necessary prerequisite to its data consistency. 

 
###CHORD RING MAINTENANCE


Every participating node is a member of a ring. The node has links to a successor node (whose key is next highest) and a predecessor (whose key is closest past). The successor is also defined as the smallest key entry stored in each routing table. Chord ring topology is defined by the successor pointer. To maintain this topology, the system is required to have exactly one ring and all members in the ring in right order. [Zave]

**Member changes:**  Before any addition or removal of nodes proceeds, the ring's member relation is redefined. When the previous example of adding a node is about to occur, C first asks the predecessor A for its successor key. The returning key is B that becomes C's  successor; C becomes B's predecessor (represented as dashed line in Figure 3).   Until A notices that B does not consider A as a  predecessor any more, B has two members that think itself as a successor. Figure 4 shows the case of C leaving. A has a missing successor; B has a missing predecessor.

<img border=0 src='figure3.png'>

> Figure 3 (left).   Adding node C.

> Figure 4 (right).   Removing node C.


<img border=0 src='figure5.png'>

> Figure 5 (left). A global view of  Figure 3.

> Figure 6 (right). A global view of Figure 4.


**Stabilization protocol:** The Chord paper suggests a periodic heartbeat feature to detect incorrect successor-relations as appeared in Figure 5 and 6. According to this, every member pings the node which have been thought to be the processor and ask for its predecessor's identitier. For example, node A asks node B for B's predecessor and B returns node C. Since the target node's predecessor is not the sender itself, A adopts B's predecessor C as new successor, knowing that C's ID is apparently greater than A's and smaller than B's.  After C sets A as the predecessor, it continues to update routing tables and key responsibilities as appropriate. As the paper does not provide any more details, in the project the member sends two arbitrary messages expecting one back, as described on the page 3 of the proposal.

An obvious weakness of this method, also mentioned in the Chord paper, is that it cannot handle concurrent joins. If more than two nodes C, D joined between two members A, B before the periodic heartbeat incorporates at least one node into the ring, B eventually sets its predecessor pointer to later added D. When A gets stabilized, it recognizes D as a successor and gives proper key responsibility to D. Although this does not violate Zave's conditions for a valid overlay structure (i.e., exactly one ring and ordered), C becomes non-reachable from other members. 

<img border=0 src='figure7.png'>

> Figure 7. Concurrent joins leaving C non-reachable.

**Successor-list:**  Note that above stabilization method cannot progress when there is 'hole in the network' as shown in Figure 6. Since the previous successor C has left the ring, A cannot get a  response. A solution from the Chord paper is that each member stores a replication of upcoming successors' keys as part of the routing table. In the example in Figure 8, current successor of node 32 is node 40. If node 40 and 52 have failed at the moment of stabilization, the successor list with replication factor 5 (the second table) can find the next live successor node 60. These nodes, whose keys are stored for the purpose of finding an alternate successor, are also called 'short links' – whereas 'long links' refers to the nodes whose keys are stored for routing efficiency. [Freedman]

<img border=0 src='figure8.png'>

> Figure 8. The Successor List.


***

I am working on it UPTO this point. (including the abstract)  -- mimi

Your proofreading would be more than welcomed :)

***

###IMPLEMENTATION

* I think it should be about how we integrate media streaming code into Chord DHT and BRBR... which I can't write about. (mimi)

* `This section should not be very long. We do not want to see your class diagrams, code snippets, etc. The point of this section is to tell us in brief notable aspects of the implementation. For example, if you have a front-end, then this is the place to say that it is written in 10K lines of JavaScript, uses WebGL for rendering robot movement, has 0 tests, and yet still seems to work. Describe libraries that you've built on, operating system dependencies, and other important details that are implementation-specific. But, please be selective.`


###EVALUATION

* the results from the streaming? (mimi)

* `How do you know that your system does what you want it to? How did you test your system? Under what scenarios/workload did it work and for how long? This is also the section where you should explain how you integrated ShiViz into your prototype -- what are events in your system, what messages are instrumented with ShiViz?`

* `If you have graphs/tables with performance information, or other measurements that you've performed on your system, then add them here. Note that all empirical evaluation results must have a proper methodology to introduce the results. What was the goal of the evaluation: why did you measure what you measured? How many nodes were in the experiment, how were they connected, did you crash a node 1min into the experiment? Typically, the more information you provide to describe the experiments, the better. But, it requires careful judgment to report just the important details.`

###LIMITATIONS

* `Every system design is predicated on reasoning about and introducing trade-offs. You can't have it all when you build complex systems (e.g., CAP theorem is a concise example of this). In this section you should explain the limitations of your system and where your system does not work (because you have tested it and it failed to work), or may not work (a hypothesis drawn from design, but that you did not test). Every design has trade-offs. Don't worry about this section being long. Focus on limitations that are fundamental, rather than incidental (e.g., doesn't work on windows because we depend on bash is an incidental limitation).`

###DISCUSSION
* `Take a step back. What worked and didn't work during the project? Tell us an interesting story about your experience in building the system. Talk about the issues and challenges that came up during the project.`

###REFERENCES

[Stoica] I. Stoica, R. Morris, D. Karger, M. F. Kaashoek, and H. Balakrishnan. Chord: A Scalable Peer-to-peer Lookup Service for Internet Applications. In Proc. of SIGCOMM. ACM, 2001.

[Karger] D. Karger, E. Lehman, T. Leighton, Matthew Lewvine, D. Lewin, and R. Panigrahy.  Consistent Hashing and Random Trees: Distributed Caching Protocols for Relieving Hot Spots on the World Wide Web. In Proc. of the 29th Annual ACM Symposium on Theory of Computing, 1997.

[Freedman] M. J. Freedman, K. Lakshminarayana, S. Rhea, and I. Stoica. Non-Transitive Connectivity and DHTs. In Proc. of WORLDS, 2005

[Zave] P. Zave.  Using lightweight modeling to understand chord. ACM CCR, 2012.

