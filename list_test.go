package proxylist

import (
	"bytes"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testfile    = "proxylist.txt"
	testfileBad = "proxylist_bad.txt"
)

var (
	ip4 = []byte(`
http://proxy1
http://proxy2
http://proxy3
http://proxy4
bad-proxy.com
`)
	ip6 = []byte(`
http://proxy1
http://proxy2
http://proxy3
http://proxy4
bad-proxy.com
`)
)

func TestList_FromReader(t *testing.T) {
	var (
		err    error
		badIP4 []string
		badIP6 []string
	)

	l := NewList()
	badIP4, err = l.FromReader(bytes.NewReader(ip4), Ip4)
	require.Nil(t, err)
	require.Len(t, badIP4, 1)
	badIP6, err = l.FromReader(bytes.NewReader(ip6), Ip6)
	require.Nil(t, err)
	require.Len(t, badIP6, 1)

	resource1 := "parser_site.com"
	resource2 := "parser_site.ru"

	_ = l.GetFree(resource1, Ip4)
	assert.Equal(t, l.NumFreeIP4(resource1), 3)
	assert.Equal(t, l.NumBusyIP4(resource1), 1)
	assert.Equal(t, l.NumFreeIP4(resource2), 4)
	assert.Equal(t, l.NumBusyIP4(resource2), 0)

	_ = l.GetFreeIP4(resource1)
	_ = l.GetFree(resource1, Ip4)
	_ = l.GetFree(resource1, Ip4)
	assert.Equal(t, l.NumFreeIP4(resource1), 0)
	assert.Equal(t, l.NumBusyIP4(resource1), 4)
	assert.Equal(t, l.NumFreeIP4(resource2), 4)
	assert.Equal(t, l.NumBusyIP4(resource2), 0)

	l.SetFree(resource1, &url.URL{
		Scheme: "http",
		Host:   "proxy2",
	}, Ip4)
	l.SetFree(resource1, &url.URL{
		Scheme: "http",
		Host:   "proxy12", // no such proxy
	}, Ip4)
	assert.Equal(t, l.NumFreeIP4(resource1), 1)
	assert.Equal(t, l.NumBusyIP4(resource1), 3)
	assert.Equal(t, l.NumFreeIP4(resource2), 4)
	assert.Equal(t, l.NumBusyIP4(resource2), 0)

	p := l.GetFree(resource1, -1) // Unknown type
	assert.Nil(t, p)

	p = l.GetFree(resource2, Ip6)
	assert.Equal(t, l.NumFreeIP4(resource1), 1)
	assert.Equal(t, l.NumBusyIP4(resource1), 3)
	assert.Equal(t, l.NumFreeIP6(resource2), 3)
	assert.Equal(t, l.NumBusyIP6(resource2), 1)
	l.SetFree(resource2, p, Ip6)
	assert.Equal(t, l.NumFreeIP4(resource1), 1)
	assert.Equal(t, l.NumBusyIP4(resource1), 3)
	assert.Equal(t, l.NumFreeIP6(resource2), 4)
	assert.Equal(t, l.NumBusyIP6(resource2), 0)

	l = NewList()

	badIP4, err = l.FromReaderIP4(bytes.NewReader(ip4))
	assert.Nil(t, err)
	assert.Len(t, badIP4, 1)

	badIP6, err = l.FromReaderIP6(bytes.NewReader(ip6))
	assert.Nil(t, err)
	assert.Len(t, badIP6, 1)
}

func TestList_FromFile_bad(t *testing.T) {
	var (
		badIP4 []string
		badIP6 []string
		err    error
	)

	p := NewList()
	badIP4, err = p.FromFile(testfileBad, Ip4)
	require.Nil(t, err)

	resource1 := "parser_site.com"
	resource2 := "parser_site.ru"

	assert.Len(t, badIP4, 2)
	assert.Equal(t, p.NumIP4(), 4)
	assert.Equal(t, p.NumFreeIP4(resource1), p.NumIP4())
	assert.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	assert.Equal(t, p.NumIP6(), 0)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	badIP6, err = p.FromFileIP6(testfileBad)
	require.Nil(t, err)

	assert.Len(t, badIP4, 2)
	assert.Equal(t, p.NumIP4(), 4)
	assert.Equal(t, p.NumFreeIP4(resource1), p.NumIP4())
	assert.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	assert.Len(t, badIP6, 2)
	assert.Equal(t, p.NumIP6(), 4)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())
}

func TestList_FromFile_good(t *testing.T) {
	var (
		badIP4 []string
		badIP6 []string
		err    error
	)

	resource1 := "parser_site.com"
	resource2 := "parser_site.ru"

	p := NewList()
	badIP4, err = p.FromFile(testfile, Ip4)
	require.Nil(t, err)

	require.Len(t, badIP4, 0)
	require.Equal(t, p.NumIP4(), 6)
	require.Equal(t, p.NumFreeIP4(resource1), p.NumIP4())
	require.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	require.Equal(t, p.NumIP6(), 0)
	require.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	require.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	badIP6, err = p.FromFile(testfile, Ip6)
	require.Nil(t, err)

	require.Len(t, badIP4, 0)
	require.Equal(t, p.NumIP4(), 6)
	require.Equal(t, p.NumFreeIP4(resource1), p.NumIP4())
	require.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	require.Len(t, badIP6, 0)
	require.Equal(t, p.NumIP6(), 6)
	require.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	require.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	proxy := p.GetFree(resource1, Ip4)
	require.NotNil(t, proxy)
	require.Equal(t, p.NumBusyIP4(resource1), 1)
	require.Equal(t, p.NumFreeIP4(resource1), 5)
	require.Equal(t, p.NumBusyIP6(resource1), 0)
	require.Equal(t, p.NumFreeIP6(resource1), 6)
	require.Equal(t, p.NumBusyIP4(resource2), 0)
	require.Equal(t, p.NumFreeIP4(resource2), 6)
	require.Equal(t, p.NumBusyIP6(resource2), 0)
	require.Equal(t, p.NumFreeIP6(resource2), 6)
}

func TestList_FromFile(t *testing.T) {
	var (
		f      *os.File
		badIP4 []string
		err    error
	)

	f, err = os.Open(testfile)
	require.Nil(t, err)

	p := NewList()
	badIP4, err = p.FromReader(f, Ip4)
	require.Nil(t, err)
	require.Len(t, badIP4, 0)
	require.Nil(t, f.Close())

	resource1 := "parser_site.com"
	resource2 := "parser_site.ru"

	assert.Equal(t, p.NumIP4(), 6)
	assert.Equal(t, p.NumFreeIP4(resource1), p.NumIP4())
	assert.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	assert.Equal(t, p.NumIP6(), 0)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	proxy := p.GetFree(resource1, Ip6)
	assert.Nil(t, proxy)
	proxy = p.GetFree(resource1, Ip4)
	assert.NotNil(t, proxy)
	assert.Equal(t, p.NumBusyIP4(resource1), 1)
	assert.Equal(t, p.NumFreeIP4(resource1), p.NumIP4()-1)
	assert.Equal(t, p.NumBusyIP4(resource2), 0)
	assert.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	assert.Equal(t, p.NumIP6(), 0)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	badIP4, err = p.FromFile(testfile, Ip4)
	assert.Nil(t, err)
	assert.Len(t, badIP4, 0)
	assert.Equal(t, p.NumIP4(), 6)
	assert.Equal(t, p.NumFreeIP4(resource1), p.NumIP4())
	assert.Equal(t, p.NumFreeIP4(resource2), p.NumIP4())
	assert.Equal(t, p.NumIP6(), 0)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	_, err = p.FromFile("not found", Ip4)
	require.NotNil(t, err)
	_, err = p.FromReader(f, Ip4)
	require.NotNil(t, err)
}

func TestList(t *testing.T) {
	resource1 := "parser_site.com"
	resource2 := "parser_site.ru"

	p := NewList()
	bad, err := p.FromFile(testfile, Ip6)
	assert.Nil(t, err)
	assert.Len(t, bad, 0)
	assert.Equal(t, p.NumIP6(), 6)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	proxy1 := p.GetFree(resource1, Ip6)
	assert.NotNil(t, proxy1)
	assert.Equal(t, p.NumBusyIP6(resource1), 1)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6()-1)
	assert.True(t, p.isBusy(resource1, proxy1, Ip6))
	assert.Equal(t, p.NumBusyIP6(resource2), 0)
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	proxy2 := p.GetFree(resource1, Ip6)
	assert.NotEqual(t, proxy1, proxy2)
	assert.Equal(t, p.NumBusyIP6(resource1), 2)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6()-2)
	assert.Equal(t, p.NumBusyIP6(resource2), 0)
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	p.SetFreeIP6(resource1, proxy1)
	assert.Equal(t, p.NumBusyIP6(resource1), 1)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6()-1)
	assert.Equal(t, p.NumBusyIP6(resource2), 0)
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	p.setBusyIP6(resource1, proxy1)
	assert.Equal(t, p.NumBusyIP6(resource1), 2)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6()-2)
	assert.Equal(t, p.NumBusyIP6(resource2), 0)
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	p.SetFree(resource1, proxy1, Ip6)
	p.SetFree(resource1, proxy2, Ip6)
	assert.Equal(t, p.NumBusyIP6(resource1), 0)
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())
	assert.Equal(t, p.NumBusyIP6(resource2), 0)
	assert.Equal(t, p.NumFreeIP6(resource2), p.NumIP6())

	// everyone is busy
	for i := 0; i < p.NumIP6()+1; i++ {
		proxy1 = p.GetFreeIP6(resource2)
		if i >= p.NumIP6() { // no free
			assert.Nil(t, proxy1)
		}
	}
	assert.Equal(t, p.NumBusyIP6(resource2), p.NumIP6())
	assert.Equal(t, p.NumFreeIP6(resource1), p.NumIP6())

	bad, err = p.FromFile(testfile, Ip6)
	assert.Nil(t, err)
	assert.Len(t, bad, 0)
	assert.Equal(t, p.NumBusyIP6(resource1), 0)
	assert.Equal(t, p.NumBusyIP6(resource2), 0)
	assert.Equal(t, p.NumIP6(), 6)
}

func TestList_Index(t *testing.T) {
	p := NewList()
	bad, err := p.FromFile(testfile, Ip6)
	require.Nil(t, err)
	require.Len(t, bad, 0)

	resource := "parser_site.com"

	proxy := p.GetFreeIP6(resource)
	assert.NotEqual(t, p.IndexIP6(proxy), -1)
	assert.Equal(t, p.IndexIP6(&url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}), -1)

	assert.Equal(t, p.IndexIP4(proxy), -1)
	assert.Equal(t, p.IndexIP4(&url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}), -1)
}

func TestList_BusyFree(t *testing.T) {
	p := NewList()
	bad, err := p.FromFileIP4(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)

	resource := "parser_site.com"

	unknownProxy := &url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}

	proxyFromTestfile := &url.URL{
		Scheme: "https",
		Host:   "12.33.12.34:12",
	}

	p.setBusyIP4(resource, proxyFromTestfile)
	p.setBusy(resource, unknownProxy, Ip4)
	assert.True(t, p.isBusy(resource, proxyFromTestfile, Ip4))

	p.SetFreeIP4(resource, unknownProxy)
	assert.False(t, p.isBusyIP4(resource, unknownProxy))
	assert.False(t, p.isBusyIP6(resource, unknownProxy))
}

func TestList_refresh(t *testing.T) {
	var (
		bad []string
		err error
	)
	p := NewList()
	_, err = p.refresh([]byte("http://127.0.0.1"), -1)
	assert.NotNil(t, err)

	bad, err = p.refresh([]byte("fdgsdf"), Ip6)
	assert.Nil(t, err)
	assert.Len(t, bad, 1)
}

func TestList_String(t *testing.T) {
	p := NewList()
	_, _ = p.FromFile(testfile, Ip6)
	_, _ = p.FromFileIP4(testfile)

	sIP4 := p.StringIP4()
	sIP6 := p.StringIP6()
	assert.Len(t, sIP4, 201)
	assert.Equal(t, len(sIP6), len(sIP4))
	assert.Equal(t, len(p.String()), len(sIP4)+len(sIP6)+1) // + \n

	p = NewList()
	_, _ = p.FromFileIP4(testfile)

	sIP4 = p.StringIP4()
	sIP6 = p.StringIP6()
	assert.Len(t, sIP4, 201)
	assert.Equal(t, len(sIP6), 57)
	assert.Equal(t, len(p.String()), len(sIP4)+len(sIP6)+1) // + \n
}
