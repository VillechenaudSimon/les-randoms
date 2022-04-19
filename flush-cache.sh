#!/bin/bash

function statAll() {
    for f in musics/*/*; do
        stat --printf="%s\t$f\n" $f
    done
}
i=0
for f in $(statAll | grep -P "^0\t" | cut -f 2); do
    rm $f
    let i=i+1
done
echo "$i files deleted in cache."
