FROM ubuntu

RUN apt-get update
RUN apt-get install -y git wget curl telnet ceph zip

RUN wget -qO- -O tmp.zip https://github.com/AcalephStorage/kviator/releases/download/v0.0.5/kviator-0.0.5-linux-amd64.zip && unzip tmp.zip && mv kviator /usr/local/bin/
