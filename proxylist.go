package proxylist

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

type List struct {
	sync.RWMutex
	items []*proxy
}

type proxy struct {
	busy bool
	url  *url.URL
}

// NewProxylist returns a structure pointer.
func NewProxylist() (l *List) {
	l = &List{
		items: []*proxy{},
	}
	return
}

// String returns a proxy list with busy tags (+/-).
func (l *List) String() string {
	l.RLock()
	defer l.RUnlock()

	res := ""
	for _, item := range l.items {
		busy := "-"
		if item.busy {
			busy = "+"
		}
		res += item.url.String() + " " + busy + "\n"
	}
	return res
}

// FromFile reads a list from a file.
func (l *List) FromFile(filename string) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	return l.parsing(b)
}

// FromReader reads a list from io.Reader.
func (l *List) FromReader(r io.Reader) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	return l.parsing(b)
}

// parsing parses proxy list.
func (l *List) parsing(b []byte) (bad []string, err error) {
	l.Lock()
	defer l.Unlock()

	bad = []string{}
	l.items = []*proxy{}
	for _, row := range bytes.Split(b, []byte("\n")) {
		row = bytes.TrimSpace(row)
		if len(row) == 0 {
			continue
		}

		u, err := url.ParseRequestURI(string(row))
		if err != nil || u.Host == "" {
			bad = append(bad, string(row))
			continue
		}

		l.items = append(l.items, &proxy{
			busy: false,
			url:  u,
		})
	}

	rand.Seed(time.Now().UnixNano())
	shuffleProxy := make([]*proxy, len(l.items))
	for i, j := range rand.Perm(len(l.items)) {
		shuffleProxy[i] = l.items[j]
	}

	l.items = shuffleProxy
	return
}

// GetFree returns a free proxy and marks it as busy.
func (l *List) GetFree() *url.URL {
	l.Lock()
	defer l.Unlock()

	// случайное прокси
	rand.Seed(time.Now().UnixNano())
	shuffleProxy := make([]*proxy, len(l.items))
	for i, j := range rand.Perm(len(l.items)) {
		shuffleProxy[i] = l.items[j]
	}

	for _, item := range shuffleProxy {
		if !item.busy {
			item.busy = true
			return item.url
		}
	}

	return nil
}

// SetFree removes the busy flag from the proxy.
func (l *List) SetFree(u *url.URL) {
	ind := l.Index(u)
	if ind < 0 {
		return
	}

	l.Lock()
	l.items[ind].busy = false
	l.Unlock()
}

// isBusy returns whether the proxy is busy.
// Always returns false if no proxy is found.
func (l *List) isBusy(u *url.URL) bool {
	ind := l.Index(u)
	if ind < 0 {
		return false
	}

	l.RLock()
	defer l.RUnlock()
	return l.items[ind].busy
}

// setBusy sets the busy flag to the proxy.
func (l *List) setBusy(u *url.URL) {
	ind := l.Index(u)
	if ind < 0 {
		return
	}

	l.Lock()
	l.items[ind].busy = true
	l.Unlock()
}

// Num returns the total number of proxies.
func (l *List) Num() int {
	l.RLock()
	defer l.RUnlock()
	return len(l.items)
}

// NumBusy returns the number of busy proxies.
func (l *List) NumBusy() (numBusy int) {
	l.RLock()
	for _, item := range l.items {
		if item.busy {
			numBusy++
		}
	}
	l.RUnlock()
	return
}

// NumFree returns the number of free proxies.
func (l *List) NumFree() (numFree int) {
	return l.Num() - l.NumBusy()
}

// Index returns the proxy index in the array.
// returns -1 if not found.
func (l *List) Index(u *url.URL) int {
	l.RLock()
	defer l.RUnlock()

	for i := range l.items {
		if l.items[i].url.String() == u.String() {
			return i
		}
	}
	return -1
}
