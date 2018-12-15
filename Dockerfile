FROM golang:1.11.3-stretch

RUN go get -u cloud.google.com/go/cmd/go-cloud-debug-agent

RUN go get -u cloud.google.com/go/profiler

RUN go get -u go.opencensus.io/trace

RUN go get -u contrib.go.opencensus.io/exporter/stackdriver

WORKDIR /go/src/hello-app
COPY main.go .
RUN go build -gcflags=all='-N -l'

COPY source-context.json .

ENV PORT 8080
ENV GOOGLE_CLOUD_PROJECT serverless-demo-219814

CMD go-cloud-debug-agent -appmodule=main -appversion=1.0 -- ./hello-app
