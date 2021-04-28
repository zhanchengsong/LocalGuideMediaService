FROM golang:1.16 AS build
WORKDIR /src
COPY . .
WORKDIR /src/server
RUN go build
EXPOSE 50051
RUN chmod +x ./server
CMD ["./server"]