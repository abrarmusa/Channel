# README.md

The project is structured into packages - each module implementing it's own
functionality. The main ones include: CustomChord, FileTransfer, StreamingServer
and StreamingClient.

To run a node you can run the following command in the src the command

`go run controller.go arg0 arg1 arg2 arg3 arg4`

where:
arg0: udp chord address for this node e.g :14321
arg1: udp chord address of a known node in the system e.g :14322
arg2: tcp rpc streaming server address e.g :1545
arg3: tcp rpc streaming client address e.g :1237
arg4: node name : has to be the same as folder in dir structure e.g node0

Note that when arg0 == arg1 that means that this is the first node to join
the system.

Instructions:
After running the above for one node, do the same for however many nodes you wish
to connect. After they're connected (you should see some finger table prints),
use node0 (since it has the sample.mp4 file we're using to test).
Type 'sample.mp4' to trigger segmenting of video into frames. After it has split
the file, the node will calculate identifiers for each segment and distribute the
frames (*.png) to respective nodes. This is how our system takes care of load
balancing. Once that is done, you can test streaming by simply typing the filename:
'sample.mp4'. You can manually look at the way the file was distributed over nodes
by looking in the directory 'sample' under the node name directory.

NOTE:
Heartbeats are overloading our nodes so we need to optimize that. Basic functionality
is in place.
