## We specify the base image we need for our
## go application
FROM golang:1.19.7-alpine3.16

# Add environment variables
# Add environment variables
ENV RAPIDAPI_PROXY_SECRET="78f5b3e0-d3d0-11ed-bf92-43930995aeef"
ENV PORT="8080"

## Install C libraries for CGO if required
RUN apk add --no-cache musl-dev

## We create an /app directory within our
## image that will hold our application source
## files
RUN mkdir /app

## We copy everything in the root directory
## into our /app directory
ADD . /app

## We specify that we now wish to execute 
## any further commands inside our /app
## directory
WORKDIR /app

## Add this go mod download command to pull in any dependencies
RUN go mod download

## Add this line to install gcc
RUN apk add --no-cache gcc

## we run go build to compile the binary
## executable of our Go program
# RUN CGO_ENABLED=1 go build -o main .
RUN go build -o main .

## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]