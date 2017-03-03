## piladb Dockerfile

This repository contains a **Dockerfile** of [piladb](https://github.com/fern4lvarez/piladb)
for [Docker](https://www.docker.com/)'s [automated build](https://registry.hub.docker.com/u/fern4lvarez/piladb/)
published to the public [Docker Hub Registry](https://registry.hub.docker.com/).
This means that the public image will always be synced with the latest source
code version.

This image has two purposes. One is to containerize the execution for piladb,
and the other to provide a portable development environment. Every container
running the `piladb` image is provided with tools like `git` or `vim.`

Also, there is a  **Makefile** that groups the most common tasks
for an easier usage of the `piladb` image.


### Base Docker Image

* [golang](https://hub.docker.com/_/golang/)


### Installation

1. Install [Docker](https://www.docker.com/).

2. Download an [automated build](https://registry.hub.docker.com/u/fern4lvarez/piladb/)
   from the public [Docker Hub Registry](https://registry.hub.docker.com/): `docker pull fern4lvarez/piladb`

   Alternatively, you can build an image from Dockerfile: `docker build -t="fern4lvarez/piladb" .`.


### Usage

#### Build image

```
make build
```

#### Push image to Docker Hub Registry

```
make push
```

#### Pull latest image from Docker Hub Registry

```
make pull
```

#### Run `pilad` server

```
make run
```

#### Start bash session into `piladb` container

```
make bash
```
