FROM golang:alpine AS build-stage

RUN apk update && apk add --no-cache git

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# BUILD THE LIBRARY
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM scratch 

COPY --from=build-stage /src/main /src/main

ENTRYPOINT ["/src/main"]
