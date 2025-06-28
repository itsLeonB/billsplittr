FROM golang:1.24-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags='-w -s' -o /billsplittr ./cmd/api/main.go

FROM gcr.io/distroless/static-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /billsplittr /billsplittr

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/billsplittr"]
