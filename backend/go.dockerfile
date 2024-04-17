FROM golang:1.18.1-alpine3.13

WORKDIR /app

# Will copy files into the /app folder as this is docker's workdir
COPY . . 

# Download and install the dependencies

RUN go get -d -v ./...


# Build the Go app and output the binary file

RUN go build -o backend_db . 

EXPOSE 8000

CMD [ "./backend_db" ]