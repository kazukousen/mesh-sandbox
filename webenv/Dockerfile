FROM golang:1.13 as base
WORKDIR /work
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webenv  -a -ldflags '-extldflags "-static"'

FROM alpine:3.9
COPY --from=base /work/webenv .
CMD [ "./webenv" ]
