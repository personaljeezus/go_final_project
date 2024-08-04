FROM golang:1.21.5

WORKDIR /go_final_project

ADD scheduler.db ./

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go_final_project

CMD ["/final"]