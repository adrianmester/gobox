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

## Future improvements

* A production-ready version of this utility would use mTLS to authenticate the
  server and the client using server and client certificates.
  
* The rsync algorithm uses a [rolling checksum](https://rsync.samba.org/tech_report/node3.html)
  which can identify identical blocks at any arbitrary
  position in the file, thus making it more efficient. For example, with the
  current implementation, adding 1 byte
  at the beginning of a large file would mean the entire file needs to be resent,
  since none of the 1KB chunks would still match. A future version of this utility
  could possibly implement an algorithm more similar to the one used by `rsync`.

* When transferring large files with non-repeating chunks (such as a several GB zip
  file), the 1KB chunks are probably inefficient. Instead of having a constant
  chunk size, it should be able to be determined for each file individually, while
  keeping the 1KB minimum chunk size.
  
  For example, `chunkSize = max(1024, fileSize/100)`

* There are some edge-cases not covered by the current tests which
  would break syncing. Identical chunks are identified by the chunkID of the first
  occurrence. If the file containing the first occurrence were to be deleted, the
  server would have no way of writing that chunk again. This can be solved by
  creating a way for the server to request chunks it doesn't have the data for.
  
  
## Protocol

`gobox` uses protobuf to communicate between the client and the server. It was
chosen because it can easily handle binary data, and the communications between
the server and client are asynchronous.

The protobuf RPC methods are defined in [`proto/gobox.proto`](./proto/gobox.proto).
