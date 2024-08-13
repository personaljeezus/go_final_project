FROM golang:1.21.5

WORKDIR /go_final

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go_final_pr

EXPOSE 7540

CMD ["/go_final_pr"]
