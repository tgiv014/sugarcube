FROM golang:1.20-buster as builder
SHELL ["/bin/bash", "-c"]

ENV NVM_DIR /usr/local/nvm
RUN mkdir -p $NVM_DIR

RUN curl --silent -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
RUN source $NVM_DIR/nvm.sh && \
    nvm install --lts=hydrogen && \
    nvm use --lts=hydrogen

ENV NODE_PATH $NVM_DIR/versions/node/v18.16.0/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v18.16.0/bin:$PATH

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN make build

FROM debian:buster-slim
VOLUME /data

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/sugarcube /app/sugarcube

ENV DB_PATH /data/db.sqlite

CMD ["/app/sugarcube"]

