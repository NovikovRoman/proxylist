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

func NewProxylist() (p *proxylist) {
	p = &proxylist{
		items: []*proxy{},
	}
	return
}

// String возвращает список прокси с метками занятости.
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

// FromFile читает список из файла
func (p *proxylist) FromFile(filename string) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	return p.load(b)
}

// FromReader читает список из io.Reader
func (p *proxylist) FromReader(r io.Reader) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	return p.load(b)
}

func (p *proxylist) load(b []byte) (bad []string, err error) {
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

// GetFree возвращает свободный прокси и помечает его как занятый.
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

// SetFree снимает флаг busy на прокси.
func (p *proxylist) SetFree(u *url.URL) {
	ind := p.Index(u)
	if ind < 0 {
		return
	}

	p.Lock()
	p.items[ind].busy = false
	p.Unlock()
}

// isBusy возвращает занят ли прокси. Всегда возвращает false, если прокси не найден.
func (p *proxylist) isBusy(u *url.URL) bool {
	ind := p.Index(u)
	if ind < 0 {
		return false
	}

	p.RLock()
	defer p.RUnlock()
	return p.items[ind].busy
}

// setBusy устанавливает флаг busy на прокси.
func (p *proxylist) setBusy(u *url.URL) {
	ind := p.Index(u)
	if ind < 0 {
		return
	}

	p.Lock()
	p.items[ind].busy = true
	p.Unlock()
}

// Num возвращает общее количество прокси.
func (p *proxylist) Num() int {
	p.RLock()
	defer p.RUnlock()
	return len(p.items)
}

// NumBusy возвращает количество занятых прокси.
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

// NumFree возвращает количество свободных прокси.
func (p *proxylist) NumFree() (numFree int) {
	return p.Num() - p.NumBusy()
}

// Index возвращает индекс прокси в массиве, если не найден возвращает -1.
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
