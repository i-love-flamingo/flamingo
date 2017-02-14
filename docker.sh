#!/bin/bash

docker build -t flamingo/akl .

docker run -ti -p 3210:3210 -v $(pwd)/akl/frontend:/go/src/flamingo/akl/frontend flamingo/akl
