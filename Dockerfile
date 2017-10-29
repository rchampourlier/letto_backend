FROM golang:1.8

WORKDIR /go/src/gitlab.com/letto/letto_backend
COPY . .

RUN git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

CMD ["go-wrapper", "run"] # ["app"]
