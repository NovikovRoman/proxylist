package proxylist

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/url"
	"os"
	"testing"
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

func TestNewList(t *testing.T) {
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

	_ = l.GetFree(Ip4)
	require.Equal(t, l.NumFreeIP4(), 3)
	require.Equal(t, l.NumBusyIP4(), 1)

	_ = l.GetFree(Ip4)
	_ = l.GetFree(Ip4)
	_ = l.GetFree(Ip4)
	require.Equal(t, l.NumFreeIP4(), 0)
	require.Equal(t, l.NumBusyIP4(), 4)

	p := l.GetFree(Ip6)
	require.Equal(t, l.NumFreeIP6(), 3)
	require.Equal(t, l.NumBusyIP6(), 1)
	l.SetFree(p, Ip6)
	require.Equal(t, l.NumFreeIP6(), 4)
	require.Equal(t, l.NumBusyIP6(), 0)
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

	require.Len(t, badIP4, 2)
	require.Equal(t, p.NumIP4(), 4)
	require.Equal(t, p.NumFreeIP4(), p.NumIP4())
	require.Equal(t, p.NumIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	badIP6, err = p.FromFile(testfileBad, Ip6)
	require.Nil(t, err)

	require.Len(t, badIP4, 2)
	require.Equal(t, p.NumIP4(), 4)
	require.Equal(t, p.NumFreeIP4(), p.NumIP4())
	require.Len(t, badIP6, 2)
	require.Equal(t, p.NumIP6(), 4)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())
}

func TestList_FromFile_good(t *testing.T) {
	var (
		/*ip4    []string
		ip6    []string*/
		badIP4 []string
		badIP6 []string
		err    error
	)

	p := NewList()
	badIP4, err = p.FromFile(testfile, Ip4)
	require.Nil(t, err)

	require.Len(t, badIP4, 0)
	require.Equal(t, p.NumIP4(), 6)
	require.Equal(t, p.NumFreeIP4(), p.NumIP4())
	require.Equal(t, p.NumIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	badIP6, err = p.FromFile(testfile, Ip6)
	require.Nil(t, err)

	require.Len(t, badIP4, 0)
	require.Equal(t, p.NumIP4(), 6)
	require.Equal(t, p.NumFreeIP4(), p.NumIP4())
	require.Len(t, badIP6, 0)
	require.Equal(t, p.NumIP6(), 6)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	proxy := p.GetFree(Ip4)
	require.NotNil(t, proxy)
	require.Equal(t, p.NumBusyIP4(), 1)
	require.Equal(t, p.NumFreeIP4(), 5)
	require.Equal(t, p.NumBusyIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), 6)

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

	require.Equal(t, p.NumIP4(), 6)
	require.Equal(t, p.NumFreeIP4(), p.NumIP4())
	require.Equal(t, p.NumIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	proxy := p.GetFree(Ip6)
	require.Nil(t, proxy)
	proxy = p.GetFree(Ip4)
	require.NotNil(t, proxy)
	require.Equal(t, p.NumBusyIP4(), 1)
	require.Equal(t, p.NumFreeIP4(), 5)
	require.Equal(t, p.NumIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	badIP4, err = p.FromFile(testfile, Ip4)
	require.Nil(t, err)
	require.Len(t, badIP4, 0)
	require.Equal(t, p.NumIP4(), 6)
	require.Equal(t, p.NumFreeIP4(), p.NumIP4())
	require.Equal(t, p.NumIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	_, err = p.FromFile("not found", Ip4)
	require.NotNil(t, err)
	_, err = p.FromReader(f, Ip4)
	require.NotNil(t, err)
}

func TestList(t *testing.T) {
	p := NewList()
	bad, err := p.FromFile(testfile, Ip6)
	require.Nil(t, err)
	require.Len(t, bad, 0)
	require.Equal(t, p.NumIP6(), 6)
	require.Equal(t, p.NumFreeIP6(), p.NumIP6())

	var (
		proxy1 *url.URL
		proxy2 *url.URL
	)

	proxy1 = p.GetFree(Ip6)
	require.NotNil(t, proxy1)
	require.Equal(t, p.NumBusyIP6(), 1)
	require.Equal(t, p.NumFreeIP6(), 5)
	require.True(t, p.isBusy(proxy1, Ip6))

	proxy2 = p.GetFree(Ip6)
	require.NotEqual(t, proxy1, proxy2)
	require.Equal(t, p.NumBusyIP6(), 2)
	require.Equal(t, p.NumFreeIP6(), 4)

	p.SetFreeIP6(proxy1)
	require.Equal(t, p.NumBusyIP6(), 1)
	require.Equal(t, p.NumFreeIP6(), 5)

	p.setBusyIP6(proxy1)
	require.Equal(t, p.NumBusyIP6(), 2)
	require.Equal(t, p.NumFreeIP6(), 4)
	p.SetFree(proxy1, Ip6)

	p.SetFree(proxy2, Ip6)
	require.Equal(t, p.NumBusyIP6(), 0)
	require.Equal(t, p.NumFreeIP6(), 6)

	// все занимаем
	for i := 0; i < p.NumIP6()+20; i++ {
		proxy1 = p.GetFreeIP6()
		if i >= p.NumIP6() { // свободных нет
			require.Nil(t, proxy1)
		}
	}

	require.Equal(t, p.NumBusyIP6(), p.NumIP6())

	bad, err = p.FromFile(testfile, Ip6)
	require.Nil(t, err)
	require.Len(t, bad, 0)
	require.Equal(t, p.NumBusyIP6(), 0)
	require.Equal(t, p.NumIP6(), 6)
}

func TestList_Index(t *testing.T) {
	p := NewList()
	bad, err := p.FromFile(testfile, Ip6)
	require.Nil(t, err)
	require.Len(t, bad, 0)

	proxy := p.GetFreeIP6()
	require.NotEqual(t, p.IndexIP6(proxy), -1)
	require.Equal(t, p.IndexIP6(&url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}), -1)

	require.Equal(t, p.IndexIP4(proxy), -1)
	require.Equal(t, p.IndexIP4(&url.URL{
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

	unknownProxy := &url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}

	p.setBusy(unknownProxy, Ip6)
	p.SetFree(unknownProxy, Ip6)
	require.False(t, p.isBusyIP6(unknownProxy))
}
