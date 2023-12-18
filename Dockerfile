FROM golang:1.21.0-alpine3.18 as base

RUN apk add --no-cache make

WORKDIR /app

COPY go.mod go.sum ./

FROM base as dev
RUN go get -u github.com/cosmtrek/air
RUN go mod download
COPY . .
RUN make build
EXPOSE 8080
CMD [ "make", "watch" ]