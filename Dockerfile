# Dockerfile Multi-Arch FreeSWITCH
FROM --platform=$BUILDPLATFORM debian:bookworm

ARG TARGETARCH
ARG TARGETOS

ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=Asia/Tehran

RUN apt update && apt install -y \
    git build-essential autoconf automake libtool pkg-config \
    libncurses5-dev libssl-dev libcurl4-openssl-dev libpcre3-dev \
    libspeexdsp-dev libedit-dev libsqlite3-dev libldns-dev libopus-dev \
    libsndfile1-dev libmpg123-dev libavformat-dev libswscale-dev wget \
    && rm -rf /var/lib/apt/lists/*

# Build FreeSWITCH
RUN mkdir -p /usr/src/freeswitch
WORKDIR /usr/src
RUN git clone https://github.com/signalwire/freeswitch.git
WORKDIR /usr/src/freeswitch
RUN ./bootstrap.sh -j && ./configure --prefix=/usr && make -j$(nproc) && make install && make samples

EXPOSE 5060/udp 5060/tcp 5061/tcp 16384-32768/udp 8021/tcp

CMD ["freeswitch", "-nonat", "-nf"]
