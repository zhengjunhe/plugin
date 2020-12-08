FROM ubuntu:16.04

WORKDIR /root
COPY dplatform dplatform
COPY dplatform-cli dplatform-cli
COPY dplatform.toml dplatform*.toml ./
COPY entrypoint.sh entrypoint.sh

CMD ["/root/dplatform", "-f", "/root/dplatform.toml"]
