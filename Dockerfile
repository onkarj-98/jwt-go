## will choose the incredebily lightweight
# Go alpine image to work with
FROM golang:1.16.1 AS builder

# will put all of our project code

RUN mkdir /app
ADD . /app
WORKDIR /app

# we want to build our application's binary executable

RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

## the lightweight scratch image will
## run our application within

FROM alpine:latest AS production

# we have to copy the output from our
# builder stage to our production stage

COPY --from=builder /app .
## we can then kickoff  our newly compiled
## binary executable
CMD ["./main"]