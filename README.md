# Hello App example

This example shows how to build and deploy a containerized Go web server
application using [Google kubernetes Engine](https://cloud.google.com/kubernetes-engine/) (a.k.a GKE).
Also this application enables Stackdriver profiler, trace and debug.

This repository contains:

- `main.go` contains the HTTP server implementation. It responds to all HTTP
  requests with a  `Hello, world!` response.
- `Dockerfile` is used to build the Docker image for the application.
- `manifests` are configuration files for deployment and service on GKE.
- `cloudbuild.yaml` is a configuration file for [Cloud Build](https://cloud.google.com/cloud-build/)