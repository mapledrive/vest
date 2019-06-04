FROM vest-main

RUN git clone --single-branch -b go https://github.com/mapledrive/vest.git /go/src/vest.com/vest

WORKDIR /go/src/vest.com/vest/go

RUN dep init

RUN dep ensure --vendor-only

# Generate binary file
RUN go build -o app server/server.go

# Expose the port
EXPOSE 8080:8080
EXPOSE 443:443

# Run the executable
CMD ["./app"]