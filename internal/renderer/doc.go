/*
Package renderer implements a PDF document generator for markdown documents.

# Introduction

This package depends on two other packages:

* The [gomarkdown](https://github.com/gomarkdown/markdown) parser to read the markdown source

* The fpdf package to generate the PDF

# Quick start

In the cmd folder is an example using the package. It demonstrates
a number of features. The test PDF was created with this command:

	go run convert.go -i test.md -o test.pdf

See README for limitations and known issues
*/
package renderer
