##Building a Chord Peer-to-Peer Overlay

###INTRODUCTION

A significant number of distributed hash table (DHT) models have been developed taking  peer-to-peer (P2P) applications explicitly into account. These prototypes include CAN, Chord [Stoica], Pastry, Kademlia and more. But it is not true that DHT exclusively serves best P2P designs. By definition, P2P systems are merely networked environment where every node is responsible both as a client and a server. DHT-based infrastructure provides the overlay network with a certain way of locating data, since it can distribute content across nodes evenly, scale up the network size and route among the nodes fast. 

The purpose of this project is to build a DHT and utilize it in a P2P application. We chose Chord among the P2P-friendly DHT systems because it has a simple design and has been frequently studied and mentioned by researchers. Chord provides accessible material for implementing load balancing, scalability and routing efficiency. However, this early protocol does not yet elaborate the issues of correctness and reliability. P2P environments “have high churn rates: nodes join and leave the overlay continuously and do not stay in the overlay for long”. [Castro] In the context of DHT alone, there are possible cases that can go wrong during the network traffic. Studies after the Chord publication have discovered and improved related correctness flaws implicated in traditional DHT designs including Chord.  Key challenge is at node failure or connection failure that the system is to maintain the topological structure.

This motivated problem is an on-going learning process and has been left as a soft requirement. As a practical goal, this project aims to perform media streaming atop the final product. It requires searching the network for multiple files and sorting chunks of media streams while tolerating inherited limitations. As a result, ….  (Summarize specific results about streaming.)

###REFERENCES

[Stoica] I. Stoica, R. Morris, D. Karger, M. F. Kaashoek, and H. Balakrishnan. Chord: A Scalable Peer-to-peer Lookup Service for Internet Applications. In Proc. of SIGCOMM. ACM, 2001.

[Castro] M. Castro, M. Costa, and A. Rowstron. Performance and Dependability of Structured Peer-to-Peer Overlays. In Proc. of DSN, 2004.
