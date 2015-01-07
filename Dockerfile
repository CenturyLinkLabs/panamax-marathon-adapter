# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/centurylinklabs/panamax-marathon-adapter

# Build dependencies
RUN go get github.com/codegangsta/martini
RUN go get github.com/jbdalido/gomarathon

# Build adapter
RUN go install github.com/centurylinklabs/panamax-marathon-adapter

# Run the adapter
ENTRYPOINT /go/bin/panamax-marathon-adapter

