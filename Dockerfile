FROM golang:1.24.0 AS build

WORKDIR /
COPY . /

RUN go build -o rbay main.go

FROM debian:bookworm-slim
RUN useradd -ms /bin/bash app
COPY --from=build --chown=app:app /rbay /home/app/
RUN chmod +x /home/app/rbay
WORKDIR /home/app
USER app

ENV PORT=8080

CMD ["/home/app/rbay"]
