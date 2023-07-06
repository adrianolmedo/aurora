FROM golang:1.20-alpine AS build

WORKDIR /src/
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go install .

ARG USERNAME
ARG USER_UID
ARG USER_GID=${USER_UID}
RUN addgroup -g ${USER_GID} -S ${USERNAME} \
    && adduser -D -u ${USER_UID} -S ${USERNAME} -s /bin/sh ${USERNAME}

FROM scratch

COPY --from=build /go/bin/aurora /bin/aurora
USER ${USERNAME}
ENTRYPOINT ["/bin/aurora"]
