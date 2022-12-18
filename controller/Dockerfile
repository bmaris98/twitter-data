FROM golang:1.19-alpine

WORKDIR /controller

RUN ls

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /controller

EXPOSE 5321

CMD [ "./controller" ]