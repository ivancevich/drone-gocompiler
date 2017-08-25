# Docker image for the docker plugin
#
#     docker build --rm=true -t ivancevich/drone-gocompiler .

FROM golang:1.9.0
RUN go get github.com/tools/godep
RUN CGO_ENABLED=0 go install -a std

ADD drone-gocompiler /bin/
ENTRYPOINT ["/bin/drone-gocompiler"]
