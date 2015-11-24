FROM golang

ENV APPDIR $GOPATH/src/github.com/netbrain/cloudfiler
RUN mkdir -p $APPDIR
ADD . $APPDIR
WORKDIR $APPDIR
RUN go get -d -v
RUN go install -v

EXPOSE 8080
ENTRYPOINT ["cloudfiler"]
