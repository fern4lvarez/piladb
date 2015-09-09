piladb [![Build Status](https://travis-ci.org/fern4lvarez/piladb.svg?branch=master)](https://travis-ci.org/fern4lvarez/piladb)
======

> _[pee-lah-dee-bee]_. _pila_ means _stack_ in Spanish.

**piladb** is a RESTful database engine based on the [stack data structure](
https://en.wikipedia.org/wiki/Stack_%28abstract_data_type%29).
Create an ilimited number of stacks, which can contain any kind of supported
data types.

Features
--------

* Stacks are auto-scalable and are only limited by the capacity of the host
  or by configuration.
* Available `pop`, `push`, `peek`, or `size` operations for each of the stacks.
* Manage databases and stacks using a REST API, so you can use it from
  any programming language.
* Totally configurable using a INI-style configuration file.
* In-memory store, with eventual persistence of data on disk.
* Written in Go, i.e. binaries are self-contained and distributable.

Development
-----------

> You need Go installed.

```bash
go get github.com/fern4lvarez/piladb/...
cd $GOPATH/src/github.com/fern4lvarez/piladb
go test -v -cover ./...
```

Credits
-------

**piladb** is documentation, test driven developed by [Fernando Álvarez](
http://www.fer.ac) on his Dell XPS 13 laptop, running Ubuntu 14.04, and
using [`vim-go`](https://github.com/fatih/vim-go) plugin within `vim` editor,
in Berlin, 2015, with the help of his dog Godín.

License
-------

MIT

