# Client library for Erlangs epmd

Some code for me to learn working with Go :) Provides an API to read the registered nodes and fetch the node details for a specific name.

## API

* `epmd.Get(nodeName) (*epmd.NodeInfo, error)`
* `epmd.Names() ([]epmd.Name, error)`

