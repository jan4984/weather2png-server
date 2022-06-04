FROM golang:1.18 as builder
WORKDIR /build

COPY go.mod ./
COPY go.sum ./
ENV GOPROXY=https://goproxy.cn
RUN go mod download
COPY . .
RUN GOARCH=arm64 GOOS=linux go build -o server cmd/server.go

FROM arm64v8/alpine:latest
WORKDIR /app
COPY --from=builder /build/resources/msyh.ttf ./
COPY --from=builder /build/server ./
ENV LOCATION_ID=101020500
ENV TTF_PATH=msyh.tff
ENV HEFENG_APIKEY="hide"
EXPOSE 10008
CMD ["./server"]
