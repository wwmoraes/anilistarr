FROM golang:1.20-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -o ./bin/ ./...


FROM scratch

COPY --from=build /src/bin/handler /usr/local/bin/handler

USER 20000:20000

CMD ["/usr/local/bin/handler"]
