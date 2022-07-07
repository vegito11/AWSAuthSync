# Start from golang base image
FROM golang:alpine as builder

# ENABLE go modules
ENV GO111MODULE=on

# Install gi 
RUN apk update && apk add --no-cache git

WORKDIR /app

# To avoid downloading depedencies every time we build image.
# we are caching all dependencies by copying go.mod and go.sum
# So if there is no change no dependencies layer won't be changed
COPY go mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# copy source code
COPY . .

# CGO_ENABLED disabled for cross system compilation
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/aws_auth_syncer .

## Start new build from scratch and use binary from last stage
FROM scratch .

COPY --from=builder /app/bin/aws_auth_syncer app/aws_auth_syncer

EXPOSE 8080

CMD [/app/aws_auth_syncer]