FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go test ./... -v
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/dooms-store.bin 

FROM scratch
COPY --from=builder /app/bin/dooms-store.bin .

EXPOSE 3000
ENTRYPOINT [ "./dooms-store.bin" ]