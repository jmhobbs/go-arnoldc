[![Build Status](https://travis-ci.org/jmhobbs/go-arnoldc.svg?branch=master)](https://travis-ci.org/jmhobbs/go-arnoldc) [![codecov](https://codecov.io/gh/jmhobbs/go-arnoldc/branch/master/graph/badge.svg)](https://codecov.io/gh/jmhobbs/go-arnoldc) [![GoDoc](https://godoc.org/github.com/jmhobbs/go-arnoldc?status.svg)](https://godoc.org/github.com/jmhobbs/go-arnoldc)

# Go ArnoldC

This is an [ArnoldC](https://github.com/lhartikk/ArnoldC) parser and interpreter written in Go using [goyacc](https://godoc.org/modernc.org/goyacc).  It was inspired by the excellent GopherCon 2018 talk, [How to Write a Parser in Go](https://www.youtube.com/watch?v=NG0s3-s3whY) by [Sugu Sougoumarane](https://twitter.com/ssougou).

The choice to implement ArnoldC was informed by Matt Steele's fantastic [GET TO THE CHOPVAR](https://www.youtube.com/watch?v=skTpd96KtiY) talk from 2015. Transpiling to JavaScript is great, but let's be honest, Go is the future ;)

I never took a compilers course because I switched from CS to MIS, so I apologize if this is horiffic.  I'm brute forcing my way through lexing, parsing and interpreter by throwing things at the wall until it works.

## Compiler

The "compiler" requires you to have a working Go install, as the source is transpiled to Go then compiled with `go build`

![Demo](https://user-images.githubusercontent.com/115059/53313796-ae22db00-3880-11e9-80a0-d00253dc216d.gif)

## Interpreter

The interpreter will parse and execute your ArnoldC program on the fly.  It works for most programs, but currently has a bug lurking somewhere that makes certain programs fail.
