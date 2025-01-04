FROM alpine:latest AS binary
ARG TARGETARCH

COPY dist /app/
RUN if [[ "${TARGETARCH}" == "arm64" ]]; then \
      mv /app/yutu_linux_arm64_v8.0/yutu-linux-arm64 /app/yutu; \
    elif [[ "${TARGETARCH}" == "amd64" ]]; then \
      mv /app/yutu_linux_amd64_v1/yutu-linux-amd64 /app/yutu; \
    fi && \
    chmod +x /app/yutu

FROM scratch AS yutu
COPY --from=binary /app/yutu /usr/local/bin/yutu
ENTRYPOINT ["/usr/local/bin/yutu"]
