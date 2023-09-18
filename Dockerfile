
FROM golang:1.21.1 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN ls -la

RUN CGO_ENABLED=0 GOOS=linux GIN_MODE=release go build cmd/mock-ehr/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage

COPY --from=build-stage /app/main /app

ENV GIN_MODE=release

EXPOSE 9999

CMD ["/app"]