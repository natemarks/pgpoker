FROM ubuntu:22.04
ARG DEBIAN_FRONTEND=noninteractive
RUN apt-get update -y && apt-get install -y git
COPY build/current/linux/amd64/pgpoker /app
ENV INTERVAL=9
CMD /app