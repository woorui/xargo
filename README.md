# xargo

[![Go](https://github.com/woorui/xargo/actions/workflows/go.yml/badge.svg)](https://github.com/woorui/xargo/actions/workflows/go.yml)

A simple implement of xargs, written in golang, Just for fun.

### Usage

```bash
$ go run main.go -h
  -C string
        command to exec (default "echo")
  -P int
        maxprocs (default 3)
  -n int
        number (default 3)
```

### Example

```bash
$ echo {1..10} | go run main.go -P 3 -n 3
7 8 9
4 5 6
1 2 3
10
```
