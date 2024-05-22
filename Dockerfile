FROM public.ecr.aws/docker/library/golang AS build
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /echo ./lambda

FROM scratch
COPY --from=build /echo /echo
ENTRYPOINT ["/echo"]
