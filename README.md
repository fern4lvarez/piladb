piladb [![Build Status](https://travis-ci.org/fern4lvarez/piladb.svg?branch=master)](https://travis-ci.org/fern4lvarez/piladb) [![GoDoc](https://godoc.org/github.com/fern4lvarez/piladb?status.svg)](https://godoc.org/github.com/fern4lvarez/piladb)
======

![Logo](http://i.imgur.com/UJaqm07.png)

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
* In-memory store, persistence on disk wannabe.
* Written in Go, i.e. binaries are self-contained and distributable.

Development
-----------

> You need Go installed.

```bash
go get github.com/fern4lvarez/piladb/...
cd $GOPATH/src/github.com/fern4lvarez/piladb
make all
```

You can also use Docker to create `piladb` builds or development environment.
Please see the [`dev`](dev/) directory.

Release
-------

> You need Docker installed.

It's possible to get `pilad` binary releases by executing `make release`.
This will cross-compile `pilad` in all available OS's and architectures.

Alternatively, if you don't have docker installed, you can release `pilad` binary
with the `make gox` command. For this, you need a configured Go environment and
[`gox`](https://github.com/mitchellh/gox) installed.

API Documentation
-----------------

Looking for the API of the Go `pila` package? [Here](https://godoc.org/github.com/fern4lvarez/piladb/pila).

Do you refer to the RESTful API? See [`pilad`](pilad/) documentation.

Credits
-------

**piladb** is documentation-, test-driven developed by [Fernando Álvarez](
http://www.fer.ac) on his Dell XPS 13 laptop, running Ubuntu 14.04, and
using [`vim-go`](https://github.com/fatih/vim-go) plugin within `vim` editor,
in Berlin and Madrid, with the support of his dog Godín.

Logo was designed by [GraphicLoads](http://www.iconarchive.com/artist/graphicloads.html).

Typography [_Lily Script One_](http://www.fontspace.com/julia-petretta/lily-script-one) designed
by [Julia Petretta](http://www.fontspace.com/julia-petretta.)

License
-------

MIT
