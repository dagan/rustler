FROM golang:1.18 AS cobra-cli
RUN go install github.com/spf13/cobra-cli@v1.3.0
WORKDIR /workspace/cmd
ENTRYPOINT ["/go/bin/cobra-cli"]

FROM golang:1.18 AS builder
WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY app app
COPY pkg pkg
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test -short ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o rustle app/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/rustle /
USER 65532:65532
ENTRYPOINT ["/rustle"]