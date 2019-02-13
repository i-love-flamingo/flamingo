#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "Rename imports"
find $PWD -type f -iname '*.go' -exec gsed -i -f $DIR/renameimports.sed "{}" +;
