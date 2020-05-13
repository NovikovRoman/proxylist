package proxylist

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"os"
	"strings"
	"testing"
)

const (
	testfile    = "proxylist.txt"
	testfileBad = "proxylist_bad.txt"
)

func TestNewProxylist(t *testing.T) {
	var (
		f   *os.File
		bad []string
		err error
	)

	p := NewProxylist()
	bad, err = p.FromFile(testfileBad)
	require.Nil(t, err)
	require.Len(t, bad, 2)
	require.Equal(t, p.Num(), 4)
	require.Equal(t, p.NumFree(), p.Num())

	bad, err = p.FromFile(testfile)
	require.Nil(t, err)
	require.Equal(t, p.Num(), 6)
	require.Equal(t, p.NumFree(), p.Num())

	proxy := p.GetFree()
	require.NotNil(t, proxy)
	require.Equal(t, p.NumBusy(), 1)
	require.Equal(t, p.NumFree(), 5)

	f, err = os.Open(testfile)
	require.Nil(t, err)
	bad, err = p.FromReader(f)
	require.Nil(t, err)
	require.Nil(t, f.Close())

	require.Equal(t, p.Num(), 6)
	require.Equal(t, p.NumFree(), p.Num())

	proxy = p.GetFree()
	require.NotNil(t, proxy)
	require.Equal(t, p.NumBusy(), 1)
	require.Equal(t, p.NumFree(), 5)

	bad, err = p.FromFile(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)
	require.Equal(t, p.Num(), 6)
	require.Equal(t, p.NumFree(), p.Num())

	_, err = p.FromFile("not found")
	require.NotNil(t, err)
	_, err = p.FromReader(f)
	require.NotNil(t, err)
}

func TestProxylist(t *testing.T) {
	p := NewProxylist()
	bad, err := p.FromFile(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)
	require.Equal(t, p.Num(), 6)
	require.Equal(t, p.NumFree(), p.Num())

	var (
		proxy1 *url.URL
		proxy2 *url.URL
	)

	proxy1 = p.GetFree()
	require.NotNil(t, proxy1)
	require.Equal(t, p.NumBusy(), 1)
	require.Equal(t, p.NumFree(), 5)
	require.True(t, p.isBusy(proxy1))

	proxy2 = p.GetFree()
	require.NotEqual(t, proxy1, proxy2)
	require.Equal(t, p.NumBusy(), 2)
	require.Equal(t, p.NumFree(), 4)

	p.SetFree(proxy1)
	require.Equal(t, p.NumBusy(), 1)
	require.Equal(t, p.NumFree(), 5)

	p.setBusy(proxy1)
	require.Equal(t, p.NumBusy(), 2)
	require.Equal(t, p.NumFree(), 4)
	p.SetFree(proxy1)

	p.SetFree(proxy2)
	require.Equal(t, p.NumBusy(), 0)
	require.Equal(t, p.NumFree(), 6)

	// все занимаем
	for i := 0; i < p.Num()+20; i++ {
		proxy1 = p.GetFree()
		if i >= p.Num() { // свободных нет
			require.Nil(t, proxy1)
		}
	}

	require.Equal(t, p.NumBusy(), p.Num())

	bad, err = p.FromFile(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)
	require.Equal(t, p.NumBusy(), 0)
	require.Equal(t, p.Num(), 6)
}

func TestProxylist_Index(t *testing.T) {
	p := NewProxylist()
	bad, err := p.FromFile(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)

	proxy := p.GetFree()
	require.NotEqual(t, p.Index(proxy), -1)
	require.Equal(t, p.Index(&url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}), -1)
}

func TestProxylist_BusyFree(t *testing.T) {
	p := NewProxylist()
	bad, err := p.FromFile(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)

	unknownProxy := &url.URL{
		Scheme: "http",
		Host:   "127.0.0.0",
		Path:   "/",
	}

	p.setBusy(unknownProxy)
	p.SetFree(unknownProxy)
	require.False(t, p.isBusy(unknownProxy))
}

func TestProxylist_String(t *testing.T) {
	p := NewProxylist()
	bad, err := p.FromFile(testfile)
	require.Nil(t, err)
	require.Len(t, bad, 0)

	p.GetFree()
	p.GetFree()

	rows := strings.Split(p.String(), "\n")
	require.Len(t, rows, p.Num()+1)

	numPlus := 0
	for i := range rows {
		if i == 6 {
			require.Equal(t, rows[i], "")
			break
		}

		require.NotEqual(t, rows[i], "")
		if rows[i][len(rows[i])-1:] == "+" {
			numPlus++
		}
	}

	require.Equal(t, numPlus, 2)
}
