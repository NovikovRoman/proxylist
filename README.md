# Proxylist

[![Build Status](https://travis-ci.com/NovikovRoman/proxylist.svg?branch=master)](https://travis-ci.com/NovikovRoman/proxylist)
[![Go Report Card](https://goreportcard.com/badge/github.com/NovikovRoman/proxylist)](https://goreportcard.com/report/github.com/NovikovRoman/proxylist)
![Codecov](https://img.shields.io/codecov/c/github/NovikovRoman/proxylist)
![GitHub](https://img.shields.io/github/license/NovikovRoman/proxylist)

> Library for working with proxy lists

## Install

```shell
go get github.com/NovikovRoman/proxylist
```

## Usage

Download proxy list:
```go
p := proxylist.NewProxylist()
bad, err := p.FromFile(testfile)
if err != nil {
    panic(err)
}
```

or

```go
p := proxylist.NewProxylist()
f, err := os.Open("proxylist.txt")
if err != nil {
    panic(err)
}

bad, err := p.FromReader(f)
if err != nil {
    panic(err)
}
```

`bad` - contains an array of invalid proxies (erroneous URL).
`err` - error reading/opening file/io.Reader.

The loaded a list shuffled.

Get a free proxy:
```go
proxy := p.GetFree()
if proxy == nil {
    panic("No free proxies.")
}
```

After using the proxy, you need to free it:

```go
p.SetFree(proxy)
```

Total number of proxies:
```go
p.Num()
```

Number of free proxies:
```go
p.NumFree()
```

Number of busy proxies:
```go
p.NumBusy()
```

## Tests

```shell
go test -race -v
```

## License
[MIT License](LICENSE) © Roman Novikov