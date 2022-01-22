FROM golang

RUN go version
ENV GOPATH=/

COPY ./ ./

# install cache
RUN apt-get update
RUN apt-get -y install redis-tools

# build go app
RUN go mod download
RUN go build -o fibonacci ./cmd/main.go

CMD ["./fibonacci"]