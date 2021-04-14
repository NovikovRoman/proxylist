# Proxylist

[![Build Status](https://travis-ci.com/NovikovRoman/proxylist.svg?branch=master)](https://travis-ci.com/NovikovRoman/proxylist)
[![Go Report Card](https://goreportcard.com/badge/github.com/NovikovRoman/proxylist)](https://goreportcard.com/report/github.com/NovikovRoman/proxylist)
![Codecov](https://img.shields.io/codecov/c/github/NovikovRoman/proxylist)
![GitHub](https://img.shields.io/github/license/NovikovRoman/proxylist)

> Library for working with proxy lists

## Install

```shell
go get github.com/NovikovRoman/proxylist/v2
```

## Usage

Download proxy list:

```go
p := proxylist.NewList()
bad, err := p.FromFile("proxylist.txt", proxylist.Ip4)
// or bad, err := p.FromFileIP4("filename")
if err != nil {
    panic(err)
}
```

or

```go
p := proxylist.NewList()
f, err := os.Open("proxylist.txt")
if err != nil {
    panic(err)
}

bad, err := p.FromReader(f, proxylist.Ip4)
// or bad, err := p.FromReaderIP4(f)
if err != nil {
    panic(err)
}
```

`bad` - contains an array of invalid proxies (erroneous URL).
`err` - error reading/opening file/io.Reader.

The loaded a list shuffled.

Get a free proxy:

```go
proxy := p.GetFreeIP4()
// or proxy := p.GetFree(proxylist.Ip4)
if proxy == nil {
    panic("No free proxies.")
}
```

After using the proxy, you need to free it:

```go
p.SetFreeIP4(proxy)
// or p.SetFree(proxylist.Ip4)
```

Total number of proxies:

```go
p.NumIP4()
// or p.Num(proxylist.Ip4)
```

Number of free proxies:

```go
p.NumFreeIP4()
// or p.NumFree(proxylist.Ip4)
```

Number of busy proxies:

```go
p.NumBusyIP4()
// or p.NumBusy(proxylist.Ip4)
```

## Tests

```shell
go test -race -v
```

## License

[MIT License](LICENSE) © Roman Novikov