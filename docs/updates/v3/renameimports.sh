#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

echo "Rename imports and usages (do my best)"
find $PWD -type f -iname '*.go' -exec gsed -i -f $DIR/renameimports.sed "{}" +;

echo "remove duplicate lines (intend is to remove duplicate imports due to consolidations)"
find $PWD -type f -iname '*.go' -exec gsed -i '$!N; /^\(.*\)\n\1$/!P; D' "{}" +;



