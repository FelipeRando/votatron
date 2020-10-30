FROM golang:1.12-alpine AS build_base

RUN apk add --no-cache git

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o votatron .

# Start fresh from a smaller image
FROM alpine:3.9 

COPY --from=build_base ./votatron ./votatron

# Run the binary program produced by `go install`
CMD ["./votatron"]
