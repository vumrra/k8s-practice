FROM golang:1.26.5-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY internal ./internal
COPY apps ./apps

ARG APP
ARG BINARY
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o "/out/${BINARY}" "./${APP}"

FROM scratch
ARG BINARY
COPY --from=build "/out/${BINARY}" /app
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/app"]
