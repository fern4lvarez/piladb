piladb
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

License
-------

MIT

