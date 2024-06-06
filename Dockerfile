FROM golang:1.22-alpine AS stage1
WORKDIR /project/ciao/

COPY go.* .
RUN  go mod download

COPY . .
RUN go build -o ./cmd/ciaoAuthSvcExec ./cmd/main.go

FROM alpine:latest
WORKDIR /project/ciao/


COPY --from=stage1 /project/ciao/cmd/ciaoAuthSvcExec ./cmd/
COPY --from=stage1 /project/ciao/dev.env ./
RUN apk update && apk add --no-cache ca-certificates
RUN apk add tzdata


EXPOSE 50051
ENTRYPOINT [ "/project/ciao/cmd/ciaoAuthSvcExec" ]