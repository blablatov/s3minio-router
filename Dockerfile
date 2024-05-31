FROM golang:1.20

RUN git clone https://github.com/blablatov/s3minio-router.git
WORKDIR s3minio-router

COPY * ./

RUN go test . && go test --bench=.

RUN CGO_ENABLED=0 GOOS=linux go build -o /s3minio-router
EXPOSE 8080

CMD ["/s3minio-router"]
