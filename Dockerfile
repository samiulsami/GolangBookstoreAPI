FROM golang:alpine AS builder
WORKDIR /src
COPY . .
RUN go build . 

FROM alpine
WORKDIR /compiled
ENV BOOKSTORE_LISTEN=0.0.0.0
COPY --from=builder /src/ /compiled/

ENTRYPOINT ["./GoBookstoreAPI"]
CMD ["serve", "--port=3000"]