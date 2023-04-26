# syntax=docker/dockerfile:1
FROM golang as go-builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY main.go ./main.go
COPY cmd ./cmd/
COPY pkg ./pkg/

RUN go mod download

RUN CGO_ENABLED=0 go build -o /tmp/eckctl /app/main.go


FROM alpine:3.17

RUN apk update --no-cache
RUN apk add bash build-base jq curl git

RUN echo "eckctl:x:1000:1000:Ecktl Non Root,,,:/home/eckctl:" >> /etc/passwd
RUN mkdir -p /home/eckctl
RUN chown eckctl: -R /home/eckctl

COPY --from=go-builder /tmp/eckctl /usr/bin/

RUN chmod 755 /usr/bin/eckctl
ENV PATH=${PATH}:/home/eckctl/.local/bin/

USER 1000

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
