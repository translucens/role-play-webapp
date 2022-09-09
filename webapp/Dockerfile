FROM golang:1.19.1-bullseye as builder

LABEL service_name="scstore"
LABEL service_role="webapp"

ENV APP_NAME=scstore
ENV ROOT=/go/src/${APP_NAME}
WORKDIR ${ROOT}

COPY . .
RUN GOOS=linux go build -v -o $APP_NAME -tags=nomsgpack


FROM gcr.io/distroless/base

LABEL service_name="scstore"
LABEL service_role="webapp"

ENV APP_NAME=scstore
ENV ROOT=/go/src/${APP_NAME}
COPY . /
COPY --from=builder ${ROOT}/${APP_NAME} ${APP_NAME}

EXPOSE 8080
ENTRYPOINT ["/scstore"]
