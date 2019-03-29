#/bin/bash

function test {
    JackCompiler.sh "$1".jack
    ./rewriting-JackCompiler "$1".jack
    d=`diff -bu "$1".vm "$1"_.vm`
}

base_uri="testcases"

for directory in `ls "$base_uri"`; do
    for file in `ls "$base_uri"/"$directory"`; do
        file_uri="$base_uri"/"$directory"/"$file"
        test "${file_uri%.*}"
    done
done
