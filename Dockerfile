FROM golang:alpine AS base
WORKDIR /gitlab
COPY . . 
RUN go build

FROM alpine:3.13
WORKDIR /gitlab
COPY --from=base /gitlab/drone-gitlab-release .

CMD [ "./drone-gitlab-release" ]
