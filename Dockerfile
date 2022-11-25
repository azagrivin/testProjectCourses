FROM golang:latest as builder

ENV GO111MODULE=on

RUN mkdir /app

COPY . /app

WORKDIR /app
RUN go generate -mod vendor ./cmd/app/ && \
    go build -mod vendor -v -o server ./cmd/app/


FROM scratch

COPY --from=builder /app/server /

CMD ["/server"]