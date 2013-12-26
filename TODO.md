TODO
====

CONFIG
------
- [ ] Dynamic config reloading
- [ ] Template config with commented out options

MAINTENANCE
-----------
- [ ] Tests
- [ ] Benchmarks
- [ ] Docs
- [ ] get/write better logging package (glog, probably)

FEATURES
--------
- [ ] Flume-style interceptors
- [ ] json -> msgpack for encoding/decoding for sqlite channel
- [ ] Working 'memory' channel (adapt ChanChannel?)
- [ ] Filesystem channel (pretty low priority)
- [ ] DB Sink (Reddis?)

MISC
----
- [ ] Headers map[string]string to map[string]interface{}? maybe map[string][]byte

BUGS
----
- [ ] filter out dummy messages on network sink

COMPLETED
=========

- [x] license (attach MIT)
- [x] github public
- [x] Vagrant
- [x] rebrand (bmx -> generic), reroot repo
- [x] rename message -> event
- [x] Generic topologies, config-based sink-channel-source bindings
  - [x] Fan out - single source, multiple channel/sinks (replication for now)
  - [x] Disconnected pipes, sets of source/channel/sinks
  - [x] fan in - multisource -> channel/sink
- [x] specify config location with command line flag