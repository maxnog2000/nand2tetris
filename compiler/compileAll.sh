#!/bin/zsh

go build -o compiler *.go
for jack in files/**/*.jack; do
	XML=${jack%%.*}.xml
	./compiler $jack > $XML
	diff=$(diff -u ${XML%/*}/xml/${XML##*/} $XML | wc -l)
	if [ $diff -ne 0 ]; then
		echo $jack;
	fi
done
rm compiler

