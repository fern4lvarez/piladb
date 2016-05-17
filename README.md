piladb [![Build Status](https://travis-ci.org/fern4lvarez/piladb.svg?branch=master)](https://travis-ci.org/fern4lvarez/piladb) [![GoDoc](https://godoc.org/github.com/fern4lvarez/piladb?status.svg)](https://godoc.org/github.com/fern4lvarez/piladb) [![osw](https://img.shields.io/badge/%E2%89%85osw-supported-blue.svg)](http://oscillating.works)
======

![Logo](http://i.imgur.com/tjQbm56.png)

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

> You need Go installed. Version 1.6+ recommended.

```bash
go get github.com/fern4lvarez/piladb/...
cd $GOPATH/src/github.com/fern4lvarez/piladb
make all
```

You can also use Docker to create `piladb` builds or development environment.
Please see the [`dev`](dev/) directory.

Dependencies
------------

**piladb** aims to minimize the amount of third party dependencies and to rely on
the Go standard library as much as possible.

Even though, it uses [`gvt`](https://github.com/FiloSottile/gvt) to vendor its few
dependencies. If you use Go 1.6+, or Go 1.5 with `GO15VENDOREXPERIMENT=1` set,
you're good to go, libraries in `vendor` directory will be used and everything will
work as expected. You don't need to install, use or understand `gvt`.

If you are using a lower Go version, **piladb** will use either your current local
version of the dependencies, or the latest one available. Things might break, so
please

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
