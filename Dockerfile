FROM golang:1.22-alpine AS stage1
WORKDIR /project/ciao/

COPY go.* .
RUN  go mod download

COPY . .
RUN go build -o ./cmd/ciaoAuthSvcExec ./cmd/main.go

FROM scratch
WORKDIR /project/ciao/


COPY --from=stage1 /project/ciao/cmd/ciaoAuthSvcExec ./cmd/
COPY --from=stage1 /project/ciao/dev.env ./

EXPOSE 3000
ENTRYPOINT [ "/project/ciao/cmd/ciaoAuthSvcExec" ]