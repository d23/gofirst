FROM debian:jessie

# Install packages.
RUN apt-get update && apt-get -y install \
  build-essential \
  golang

RUN apt-get -y install procps net-tools vim

VOLUME /opt/gofirst/
WORKDIR /opt/gofirst/

CMD ["/usr/bin/go", "build"]
