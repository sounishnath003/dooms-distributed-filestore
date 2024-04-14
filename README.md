# Dooms Distributed Filestore

Dooms Filestore - is an implementation of distributed filestore system. Which works and supports multiple `Transport` networking layers (OSI-5) methods like `TCP,UDP,Sockets,GRPC,...) currently.

## Project structure

- `transport` - contains all the library implementation of the dooms transport requirements.

## Run locally

```bash
go mod tidy
make test
make run
```
