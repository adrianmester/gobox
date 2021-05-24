# gobox

A file sync utility, similar to `rsync`, written in go.

## Algorithm

`gobox` uses a simplified version of the `rsync` algorithm. All files are split
into 1KB chunks. A map of the md5 checksum of each transferred chunk is kept in
memory on the client side in the form of `map[md5sum]ChunkID`, where `ChunkID`
is made up of a `(FileID,ChunkNumber)` tuple.

The checksum of each chunk is computed before being sent, and if it has already
been transmitted, the ChunkID of the first occurance is sent instead, without
the data. The server will then read the previous chunk from disk and write that
to th new file instead.

## Limitations

* The [original rsync algorithm](https://rsync.samba.org/tech_report/node3.html)
uses a rolling checksum which can identify identical blocks at any arbitrary
position in the file, thus making it more efficient. For example, adding 1 byte
at the beginning of a large file would mean the entire file needs to be resent,
since none of the 1KB chunks would still match. A future version of this utility
could possibly implement this functionality.

* When transferring large files with non-repeating chunks (such as several GB zip
file), the 1KB chunks are probably inefficient. Instead of having a constant
chunk size, it should be able to be determined for each file individually, while
keeping a 1KB minimum.

* There are some edge-cases not covered by the current tests which
would break syncing. Identical chunks are identified by the chunkID of the first
occurrence. If the file containing the first occurrence were to be deleted, the
server would have no way of writing that chunk again. This can be solved by
creating a way for the server to request chunks it doesn't have the data for.
  
## Protocol

`gobox` uses protobuf to communicate between the client and the server. It was
chosen because it can easily handle binary data, and the communications between
the server and client can be asynchronous.
