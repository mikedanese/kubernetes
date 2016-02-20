FROM ubuntu:wily

RUN apt-get update -q \
    && apt-get install -y ohai ansible jq curl parallel \
    && apt-get clean \
    && apt-get autoremove \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
