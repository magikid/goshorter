# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.18.3 as builder

ENV GOPROXY http://proxy.golang.org

RUN mkdir -p /src/github.com/magikid/goshorter
WORKDIR /src/github.com/magikid/goshorter

# this will cache the npm install step, unless package.json changes
COPY package.json yarn.lock ./
RUN yarn install
# Copy the Go Modules manifests
COPY go.mod go.sum ./
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . .
RUN yarn install && buffalo build --static -o /bin/app

FROM alpine:3
RUN apk add --no-cache bash==5.1.16-r2 ca-certificates==20220614-r0

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD ["exec", "/bin/app"]
