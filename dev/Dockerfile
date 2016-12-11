# piladb Dockerfile
#
# https://github.com/fern4lvarez/piladbdb
#

# Pull base image.
FROM golang:latest

# Install basic development tools.
RUN \
  apt-get update && \
  apt-get install -y git vim

# Install dependencies.
RUN \
  go get github.com/golang/lint/golint

# Install piladb.
RUN \
  go get github.com/fern4lvarez/piladb && \
  cd /go/src/github.com/fern4lvarez/piladb && \
  make all

# Define mountable directories.
VOLUME ["/go/src/github.com/fern4lvarez/piladb"]

# Define working directory.
WORKDIR /go/src/github.com/fern4lvarez/piladb

# Define default command.
CMD ["/go/bin/pilad"]

# Expose port.
EXPOSE 1205
