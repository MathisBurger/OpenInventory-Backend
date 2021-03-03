FROM golang:1.15

WORKDIR /build

COPY . .
RUN go mod tidy

RUN go build -o ./backend

RUN chmod +x ./backend

EXPOSE 8080

CMD ./backend