# FILE SHARING PACKAGE

## TO TEST
```go run downloader.go :3000 :4000```
```go run downloader.go :4000 :3000```

### 1.

The nodes will pre-process the video into byte segments(TODO: 8 byte segment size). This will take about 45 seconds. Once processed, the node will open up an RPC connection to process incoming requests. 

### 2.
From the other node you can now use the following commands to get details of files available on the node. The commands are:

1. get [filename.ext] <-- input the filename with the extension.
2. the address of the node eg. [:____] some local port no.

## CURRENTY ONLY METHOD COMPLETED. OBTAINING FILE SEGMENTS WILL BE A VERY SIMILAR METHOD.