#!/bin/bash

curl -v 172.24.160.1:8800/objests/test2 -XPUT -d "this is object test2"

curl 172.24.160.1:8800/locate/test2
echo
curl 172.24.160.1:8800/objects/test2