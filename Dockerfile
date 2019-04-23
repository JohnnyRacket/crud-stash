FROM golang:latest
ADD . /go/src/crud-stash
WORKDIR /go/src/crud-stash
RUN go get 
RUN go build
ENTRYPOINT [ "./crud-stash" ]