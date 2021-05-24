# gobox
![tests](https://github.com/adrianmester/gobox/actions/workflows/go.yml/badge.svg)

A file sync utility, similar to `rsync`, written in go.

## Algorithm

`gobox` uses a simplified version of the `rsync` algorithm. All files are split
into 1KB chunks. A map of the md5 checksum of each transferred chunk is kept in
memory on the client side in the form of a `map[md5sum]ChunkID`, where `ChunkID`
is made up of a `(FileID,ChunkNumber)` tuple.

The checksum of each chunk is computed before being sent, and if it has already
been transmitted, the ChunkID of the first occurance is sent instead, without
the contents of the chunk. The server will then read the previous chunk from disk
and write that to the new file instead.

### Flow

At startup, the client recursiveley walks the watched directory, and calls the `SendFileInfo`
RPC method for each file and folder it encounters. This sends the file name, its
type, size and modification time. The servers checks that against the files it
has on disk, and if there's a mismatch, it replies with `SendChunkIds=true`.

The client will send the Chunks for all the files the server has requested. In
the case of chunks whose md5hash was already known, just the metadata is sent
(FileID, ChunkNumber), and the server retrieves the bytes from previously sent
files.

Once all the initial files have been sent, the client calls `InitialSyncComplete`.
This tells the server that any files that are on disk that haven't been sent by
the client need to be deleted.

Next, the client uses the [`fsnotify`](https://github.com/fsnotify/fsnotify)
library to monitor for file system events on the watched directory. This library
was chosen because of its cross-platform support.

Any event in the watched directory will trigger a `SendFileInfo` method call, the
rest of the flow being the same as before. The only difference is the `DeleteFile`
method which will be called when files are deleted.

## Protocol

`gobox` uses protobuf to communicate between the client and the server. It was
chosen because it can easily handle binary data, and asynchronous streams.

The protobuf RPC methods are defined in [`proto/gobox.proto`](./proto/gobox.proto).

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
  
* rsync also supports a `--checksum` mode where the checksum of the files are verified
  on both the client and server side. This could also be added to `gobox`, although
  the file size and modification times should be enough in the vast majority of cases
  to determine if a file needs to be uploaded.
  
* In addition to the current integration tests, unit tests should be written where
  appropriate.
