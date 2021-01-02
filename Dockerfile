FROM golang:alpine

WORKDIR /dezhban 
ADD . /dezhban/
RUN go build -o /usr/bin/dezhban main.go
RUN chmod +x /usr/bin/dezhban

CMD ["/usr/bin/dezhban"]
