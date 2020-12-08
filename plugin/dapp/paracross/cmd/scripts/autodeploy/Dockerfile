FROM ubuntu:16.04

WORKDIR /root
COPY dplatform dplatform
COPY dplatform-cli dplatform-cli
COPY entrypoint.sh entrypoint.sh
COPY dplatform.toml dplatform*.toml ./


CMD ["/root/dplatform", "-f", "/root/dplatform.toml"]
