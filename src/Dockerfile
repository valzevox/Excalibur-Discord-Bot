# For deployment on render.com and others
FROM golang:latest

WORKDIR /

COPY . .

RUN go build Excalibur

EXPOSE 8080

CMD [ "./Excalibur" ]
