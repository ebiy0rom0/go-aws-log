FROM golang:1.18-buster as builder

WORKDIR /

COPY go.* ./
RUN go mod download

COPY . ./
RUN make build


FROM debian:buster-slim as prod

WORKDIR /app

COPY --from=builder serverd .

EXPOSE 1323
CMD [ "/app/serverd" ]