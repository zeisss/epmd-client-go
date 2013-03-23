# Client library for Erlangs epmd

Some code for me to learn working with Go :) Provides an API to read the registered nodes and fetch the node details for a specific name.

## API

* `epmd.Get(nodeName) (*epmd.NodeInfo, error)`
* `epmd.Names() ([]epmd.Name, error)`
* `epmd.Register(...)` - Anounce itself as a 

## Links

* [Erlang distribution protocol](http://www.erlang.org/doc/apps/erts/erl_dist_protocol.html)
* [epmd source code](https://github.com/erlang/otp/blob/maint/erts/epmd/src/)
* [eclus - go implementation of erlangd epmd](https://github.com/goerlang/eclus)


