# To build : docker build . -t verloop
# TO run : docker run --rm -p 9090:9090 --env VERLOOP_DEBUG=$VERLOOP_DEBUG -env VERLOOP_DSN=$VERLOOP_DSN -it verloop
FROM bash:4.4
WORKDIR /
COPY verloop /
COPY config.yaml /
CMD ["/verloop"]