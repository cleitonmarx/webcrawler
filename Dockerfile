#
# API Frontend Dockerfile
#

# Pull the base image
FROM golang:1.5.2
MAINTAINER Cleiton Marques <cleitonmarx@hotmail.com>

# Set GOPATH
ENV GOPATH /go

# Set App_Env
ENV APP_ENV DOCKER

# Make directories for webcrawler
RUN mkdir -p /go/src/github.com/cleitonmarx/webcrawler

# Add webcrawler files
ADD . /go/src/github.com/cleitonmarx/webcrawler

# Define working directory
WORKDIR /go/src/github.com/cleitonmarx/webcrawler

# Restore Dependencies and Install Application
RUN cd /go/src/github.com/cleitonmarx/webcrawler
RUN \
	go get github.com/tools/godep && \
	godep restore && \
	godep go install

# Define default command
CMD ["/go/bin/webcrawler"]

# Expose Ports
EXPOSE 3333
