# xargo

[![Go](https://github.com/woorui/xargo/actions/workflows/go.yml/badge.svg)](https://github.com/woorui/xargo/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/woorui/xargo/branch/master/graph/badge.svg?token=4JFSD1YEE5)](https://codecov.io/gh/woorui/xargo)

A simple implement of xargs, written in golang, Just for fun.

### Install

```bash
go get -u github.com/woorui/xargo
```

### Usage

```bash
$ xargo -h
Usage of xargo:
  -C string
        command to exec (default "echo")
  -P int
        maxprocs (default 3)
  -n int
        number (default 3)
```

### Example

```bash
$ echo {1..10} | xargo -P 3 -n 3
7 8 9
4 5 6
1 2 3
10
```
