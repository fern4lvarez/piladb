# piladb Dockerfile
#
# https://github.com/fern4lvarez/piladbdb
#

# Pull base image.
FROM golang:1.5.1

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
