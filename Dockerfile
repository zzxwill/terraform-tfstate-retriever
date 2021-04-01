FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go build -o ./terraform-tfstate-retriever .

CMD ["./terraform-tfstate-retriever"]