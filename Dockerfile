FROM golang:1.21-bullseye as build

RUN mkdir -p /git-copy/
COPY . /git-copy/
WORKDIR /git-copy

ENV GO111MODULE=on
RUN make install
RUN make build

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /git-copy/cmd/cmd /cmd

ENTRYPOINT ["/cmd"]