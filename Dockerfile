FROM golang
ADD . /go/src/servertime
RUN go install servertime
ENTRYPOINT /go/bin/servertime
EXPOSE 8080
