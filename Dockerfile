FROM golang:alpine as builder

LABEL maintainer="Vanessasaurus <@vsoch>"

WORKDIR /code

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Additional dependencies that are nice to have
RUN apk add --no-cache gcc build-base linux-headers

# Build the Go app
RUN make

# Start again with minimal envoirnment.
FROM alpine

RUN apk add --no-cache git
WORKDIR /code
COPY --from=builder /code/codestats /code/codestats
COPY --from=builder /code/entrypoint.sh /code/entrypoint.sh
ENV PATH=/code:$PATH

# Command to run the executable
ENTRYPOINT ["/bin/bash", "entrypoint.sh"]
