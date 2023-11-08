# Proxylist

[![Build Status](https://app.travis-ci.com/NovikovRoman/proxylist.svg?branch=master)](https://app.travis-ci.com/NovikovRoman/proxylist)
[![Go Report Card](https://goreportcard.com/badge/github.com/NovikovRoman/proxylist/v3)](https://goreportcard.com/report/github.com/NovikovRoman/proxylist/v3)
![Codecov](https://img.shields.io/codecov/c/github/NovikovRoman/proxylist)
![GitHub](https://img.shields.io/github/license/NovikovRoman/proxylist)

> Library for working with proxy lists

## Install

```shell
go get github.com/NovikovRoman/proxylist/v3
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
resource := "site.com" // any resource name using a proxy
proxy := p.GetFreeIP4(resource)
// or proxy := p.GetFree(resource, proxylist.Ip4)
if proxy == nil {
    panic("No free proxies.")
}
```

After using the proxy, you need to free it:

```go
resource := "site.com" // any resource name using a proxy
p.SetFreeIP4(resource, proxy)
// or p.SetFree(resource, proxylist.Ip4)
```

Total number of proxies:

```go
p.NumIP4()
// or p.Num(proxylist.Ip4)
```

Number of free proxies:

```go
resource := "site.com" // any resource name using a proxy
p.NumFreeIP4(resource)
// or p.NumFree(resource, proxylist.Ip4)
```

Number of busy proxies:

```go
resource := "site.com" // any resource name using a proxy
p.NumBusyIP4(resource)
// or p.NumBusy(resource, proxylist.Ip4)
```

## Tests

```shell
go test -race -v
```

## License

[MIT License](LICENSE) © Roman Novikov
