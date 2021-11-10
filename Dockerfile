FROM golang:1.16.10-alpine3.14 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY database ./database/
COPY models ./models/
COPY routes ./routes/
COPY main.go ./

RUN go build -o /server ./main.go


FROM alpine:latest

COPY --from=build /server /server

EXPOSE 80

ENTRYPOINT [ "/server" ]
