FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/promotion ./cmd/promotion

FROM alpine:3.20 AS final
WORKDIR /app
COPY --from=build /app/promotion .
COPY migrations/promotion ./migrations/promotion
EXPOSE 7003 7013
ENTRYPOINT ["./promotion"]
