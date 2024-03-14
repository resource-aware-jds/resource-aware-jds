# Stage 1 - Builder: Import the golang container.
FROM golang:1.21-alpine as builder

# Install ssh client and git
RUN apk add --no-cache openssh-client git

# Download public key for github.com
RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

# Set the work directory.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod ./
COPY go.sum ./

# Install the dependencies.
RUN git config --global --add url."ssh://git@github.com/".insteadOf "https://github.com/"
# This command will have access to the forwarded agent (if one is
# available)
RUN go mod download

# Copy the source code into the container.
COPY ./ ./

# Build the source code
RUN go build -o ./out/cp ./cmd/controlplane/main.go
RUN go build -o ./out/worker ./cmd/worker/main.go


# Stage 2 - Runner.
FROM docker:25.0.3-dind-alpine3.19 as runner
WORKDIR /app
COPY --from=builder /app/out/cp cp
COPY --from=builder /app/out/worker worker
COPY --from=builder /app/docker-entrypoint.sh docker-entrypoint.sh

RUN apk add openrc
RUN ln -s /usr/local/bin/dockerd /etc/init.d/dockerd
RUN ln -s /usr/local/bin/docker /usr/bin/docker
RUN rc-update add dockerd


ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD [ "/app/." ]
