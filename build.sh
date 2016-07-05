#!/bin/bash
export GOPATH=`pwd`
echo "[-] Get Packages"
go get github.com/garyburd/redigo/redis
go get github.com/julienschmidt/httprouter
go get github.com/manucorporat/sse
go get github.com/rs/cors
go get github.com/sudhirj/strobe
echo "[-] Go Build"
go build -o satellite
