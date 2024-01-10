FROM golang:1.21.6-alpine3.19

WORKDIR /usr/src/app

# install bash
RUN apk add bash

# install nodemon
RUN apk add --update nodejs npm
RUN npm update npm -g
RUN npm install nodemon -g

# Install Go dependencies
COPY ./../../src/examples/golang/go.mod .
COPY ./../../src/examples/golang/go.sum .
RUN go mod download


ENTRYPOINT (sed -i 's/\r$//' ./run.sh) ; (chmod +x ./run.sh) ; ./run.sh