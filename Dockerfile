FROM golang:1.22.2-alpine AS build
RUN apk add --no-cache gcc libc-dev
WORKDIR /go/src/app

COPY . .
RUN go test  ./...
RUN go build -mod vendor -o /bin/template-wh


FROM alpine:3.19.1
MAINTAINER Peter Reisinger <p.reisinger@gmail.com>
RUN apk add --no-cache ca-certificates

COPY --from=build /bin/template-wh /usr/local/bin/template-wh
CMD ["template-wh"]
