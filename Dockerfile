ARG ARCH="amd64"
FROM ${ARCH}/alpine

ARG OS="linux"
ARG BIN_ARCH="amd64"
ENV HOME="/app"

LABEL name="idp-scim-sync"

RUN apk add --no-cache --update \
    ca-certificates \
    && rm -rf /tmp/* /var/tmp/* /var/cache/apk/*

RUN mkdir -p $HOME && \
    chown -R nobody.nobody $HOME

COPY dist/idpscim-${OS}-${BIN_ARCH} ${HOME}/idpscim
COPY dist/idpscimcli-${OS}-${BIN_ARCH} ${HOME}/idpscimcli

ENV PATH="${PATH}:${HOME}"

VOLUME ${HOME}
USER nobody:nobody
WORKDIR $HOME

CMD ["/app/idpscim", "--help"]