ARG ARCH
FROM ${ARCH}/alpine

ARG OS="linux"
ARG BIN_ARCH="amd64"

ENV PROJECT_NAME="idp-scim-sync"
ENV HOME="/app"

LABEL name="${PROJECT_NAME}" \
  org.opencontainers.image.source="https://github.com/slashdevops/idp-scim-sync"

RUN apk add --no-cache --update \
  ca-certificates \
  && rm -rf /tmp/* /var/tmp/* /var/cache/apk/*

RUN mkdir -p $HOME && \
  chown -R nobody.nobody $HOME

COPY dist/$PROJECT_NAME-$OS-$BIN_ARCH/* $HOME/

ENV PATH="${PATH}:${HOME}"

VOLUME $HOME
USER nobody:nobody
WORKDIR $HOME

CMD ["/app/idpscim", "--help"]