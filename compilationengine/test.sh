#/bin/bash

JackCompiler.sh Script.jack
d=`diff -bu Script.vm Script_.vm`
if [ "$d" == "" ]; then
    echo OK
else
    echo FAIL
    echo "$d"
    exit 1
fi
