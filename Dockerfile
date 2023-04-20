FROM golang:alpine as build_container
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .


FROM alpine
WORKDIR /root/
COPY --from=build_container /app/server .

EXPOSE 7887

ENTRYPOINT ["./server"]
