FROM golang:alpine as builder
ARG commit
ARG commitDate
ARG version

ENV MOD="github.com/eat-pray-ai/yutu/cmd"
WORKDIR /app
COPY . .

RUN Version="${MOD}.Version=${version}" \
    Commit="${MOD}.Commit=${commit}" \
    CommitDate="${MOD}.CommitDate=${commitDate}" \
    Os="${MOD}.Os=linux" \
    Arch="${MOD}.Arch=$(go env GOARCH)" \
    ldflags="-s -X ${Version} -X ${Commit} -X ${CommitDate} -X ${Os} -X ${Arch}" \
    && export ldflags \
    && go mod download \
    && go build -ldflags "${ldflags}" -o yutu .

FROM alpine:latest as yutu

COPY --from=builder /app/yutu /usr/local/bin/yutu
RUN chmod +x /usr/local/bin/yutu

ENTRYPOINT ["/usr/local/bin/yutu"]
