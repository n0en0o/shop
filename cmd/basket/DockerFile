FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/basket ./cmd/basket

FROM alpine:3.20 AS final
WORKDIR /app
COPY --from=build /app/basket .
COPY migrations/basket ./migrations/basket
EXPOSE 7002
ENTRYPOINT ["./basket"]
