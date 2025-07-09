ARG GO_VERSION

FROM public.ecr.aws/docker/library/golang:${GO_VERSION} AS builder
ENV APPDIR /go/src/github.com/acu4git/gimme-scholarship
COPY . $APPDIR/
WORKDIR $APPDIR
ADD https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o gimme-scholarship-task cmd/task/main.go

FROM gcr.io/distroless/static
COPY --from=builder /go/src/github.com/acu4git/gimme-scholarship/gimme-scholarship-task ./
COPY --from=builder /go/src/github.com/acu4git/gimme-scholarship/global-bundle.pem ./
ENTRYPOINT ["./gimme-scholarship-task"]
