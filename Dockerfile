FROM golang:1.16.5-alpine AS build
RUN apk add --no-cache gcc libc-dev
WORKDIR /go/src/app

COPY . .
RUN go test  ./...
RUN go build -mod vendor -o /bin/template-wh


FROM alpine:3.13.5
MAINTAINER Peter Reisinger <p.reisinger@gmail.com>
RUN apk add --no-cache ca-certificates

COPY --from=build /bin/template-wh /usr/local/bin/template-wh
CMD ["template-wh"]
