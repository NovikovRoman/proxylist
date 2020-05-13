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

type proxylist struct {
	sync.RWMutex
	items []*proxy
}

type proxy struct {
	busy bool
	url  *url.URL
}

// NewProxylist returns a structure pointer.
func NewProxylist() (p *proxylist) {
	p = &proxylist{
		items: []*proxy{},
	}
	return
}

// String returns a proxy list with busy tags (+/-).
func (p *proxylist) String() string {
	p.RLock()
	defer p.RUnlock()

	res := ""
	for _, item := range p.items {
		busy := "-"
		if item.busy {
			busy = "+"
		}
		res += item.url.String() + " " + busy + "\n"
	}
	return res
}

// FromFile reads a list from a file.
func (p *proxylist) FromFile(filename string) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	return p.parsing(b)
}

// FromReader reads a list from io.Reader.
func (p *proxylist) FromReader(r io.Reader) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	return p.parsing(b)
}

// parsing parses proxy list.
func (p *proxylist) parsing(b []byte) (bad []string, err error) {
	p.Lock()
	defer p.Unlock()

	bad = []string{}
	p.items = []*proxy{}
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

		p.items = append(p.items, &proxy{
			busy: false,
			url:  u,
		})
	}

	rand.Seed(time.Now().UnixNano())
	shuffleProxy := make([]*proxy, len(p.items))
	for i, j := range rand.Perm(len(p.items)) {
		shuffleProxy[i] = p.items[j]
	}

	p.items = shuffleProxy
	return
}

// GetFree returns a free proxy and marks it as busy.
func (p *proxylist) GetFree() *url.URL {
	p.Lock()
	defer p.Unlock()

	// случайное прокси
	rand.Seed(time.Now().UnixNano())
	shuffleProxy := make([]*proxy, len(p.items))
	for i, j := range rand.Perm(len(p.items)) {
		shuffleProxy[i] = p.items[j]
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
func (p *proxylist) SetFree(u *url.URL) {
	ind := p.Index(u)
	if ind < 0 {
		return
	}

	p.Lock()
	p.items[ind].busy = false
	p.Unlock()
}

// isBusy returns whether the proxy is busy.
// Always returns false if no proxy is found.
func (p *proxylist) isBusy(u *url.URL) bool {
	ind := p.Index(u)
	if ind < 0 {
		return false
	}

	p.RLock()
	defer p.RUnlock()
	return p.items[ind].busy
}

// setBusy sets the busy flag to the proxy.
func (p *proxylist) setBusy(u *url.URL) {
	ind := p.Index(u)
	if ind < 0 {
		return
	}

	p.Lock()
	p.items[ind].busy = true
	p.Unlock()
}

// Num returns the total number of proxies.
func (p *proxylist) Num() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.items)
}

// NumBusy returns the number of busy proxies.
func (p *proxylist) NumBusy() (numBusy int) {
	p.RLock()
	for _, item := range p.items {
		if item.busy {
			numBusy++
		}
	}
	p.RUnlock()
	return
}

// NumFree returns the number of free proxies.
func (p *proxylist) NumFree() (numFree int) {
	return p.Num() - p.NumBusy()
}

// Index returns the proxy index in the array.
// returns -1 if not found.
func (p *proxylist) Index(u *url.URL) int {
	p.RLock()
	defer p.RUnlock()

	for i := range p.items {
		if p.items[i].url.String() == u.String() {
			return i
		}
	}
	return -1
}
