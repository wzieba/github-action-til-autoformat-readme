FROM golang:1.14 as builder
WORKDIR /go/src/app
COPY . .
RUN go build -mod=vendor -o /go/bin/main .

FROM alpine:latest
WORKDIR /root
RUN apk update && apk add git

COPY --from=builder /go/bin/main ./main
COPY --from=builder /go/src/app/README.md.tmpl ./README.md.tmpl
COPY --from=builder /go/src/app/README_TABLE.md.tmpl ./README_TABLE.md.tmpl
COPY --from=builder /go/src/app/entrypoint.sh ./entrypoint.sh

ENTRYPOINT [ "/root/entrypoint.sh" ]

ENV REPO_PATH "/github/workspace"
