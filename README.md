# CloudFiler
A simple web based access controlled file archive. Upload files and share them with users on a per user or group level.

## Prerequisites

* Golang (http://golang.org/doc/install)

## Download

Download the application by any of these means:

* go : `go get github.com/netbrain/cloudfiler`
* git: `git clone https://github.com/netbrain/cloudfiler`
* zip: https://github.com/netbrain/cloudfiler/archive/master.zip

## Compile

```
cd $GOPATH/src/github.com/netbrain/cloudfiler # or whatever path you downloaded to.
go build

```

## Run and configure

`./cloudfiler`

The application has now started on http://127.0.0.1:8080 by default. To make any configuration changes, 
edit the file: `$HOME/.cloudfiler/config` and restart cloudfiler.

When you access the application for the first time you are prompted by an initialization process which
creates the first admin user of the system.
