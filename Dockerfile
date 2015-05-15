FROM debian:jessie

# Install packages.
RUN apt-get update && apt-get -y install \
  build-essential procps net-tools vim git \
  golang

RUN mkdir /opt/gopath
ENV GOPATH=/opt/gopath

VOLUME /opt/gofirst/
WORKDIR /opt/gofirst/
CMD ["./build.sh"]
