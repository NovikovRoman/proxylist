package proxylist

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"math/rand"
	"net/url"
	"sync"
	"time"
)

const (
	Ip4 = iota
	Ip6
)

// List structure of proxy lists.
type List struct {
	ip4Mutex sync.RWMutex
	ip4      []*proxy

	ip6Mutex sync.RWMutex
	ip6      []*proxy
}

// NewList returns a pointer to a List structure.
func NewList() (l *List) {
	l = &List{
		ip4: []*proxy{},
		ip6: []*proxy{},
	}
	return
}

// String returns a proxy list with busy tags (+/-).
func (l *List) String() string {
	return l.StringIP4() + "\n" + l.StringIP6()
}

func (l *List) StringIP4() (s string) {
	l.ip4Mutex.RLock()
	s = listString(l.ip4, Ip4)
	l.ip4Mutex.RUnlock()
	return
}

func (l *List) StringIP6() (s string) {
	l.ip6Mutex.RLock()
	s = listString(l.ip6, Ip6)
	l.ip6Mutex.RUnlock()
	return
}

func listString(list []*proxy, typeIP int) (s string) {
	switch typeIP {
	case Ip4:
		s = "IP4"

	case Ip6:
		s = "IP6"
	}

	s = "-----------------\n" + s + " PROXY\n-----------------\n"
	if len(list) == 0 {
		s += "empty list\n"
		return
	}

	for _, item := range list {
		s += item.String() + "\n"
	}
	return
}

// FromFile reads a list from a file. typeIP - proxy type (Ip4 or Ip6).
func (l *List) FromFile(filename string, typeIP int) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	return l.refresh(b, typeIP)
}

// FromFileIP4 reads proxy IP4 list from a file.
func (l *List) FromFileIP4(filename string) (bad []string, err error) {
	return l.FromFile(filename, Ip4)
}

// FromFileIP6 reads proxy IP6 list from a file.
func (l *List) FromFileIP6(filename string) (bad []string, err error) {
	return l.FromFile(filename, Ip6)
}

// FromReader reads a list from io.Reader. typeIP - proxy type (Ip4 or Ip6).
func (l *List) FromReader(r io.Reader, typeIP int) (bad []string, err error) {
	var b []byte
	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	return l.refresh(b, typeIP)
}

// FromReaderIP4 reads proxy IP4 list from io.Reader.
func (l *List) FromReaderIP4(r io.Reader) (bad []string, err error) {
	return l.FromReader(r, Ip4)
}

// FromReaderIP6 reads proxy IP6 list from io.Reader.
func (l *List) FromReaderIP6(r io.Reader) (bad []string, err error) {
	return l.FromReader(r, Ip6)
}

// GetFree returns a free proxy and marks it as busy. typeIP - proxy type (Ip4 or Ip6).
func (l *List) GetFree(resource string, typeIP int) *url.URL {
	var shuffleProxy []*proxy

	switch typeIP {
	case Ip4:
		l.ip4Mutex.Lock()
		defer l.ip4Mutex.Unlock()

		// случайное прокси
		rand.Seed(time.Now().UnixNano())
		shuffleProxy = make([]*proxy, len(l.ip4))
		for i, j := range rand.Perm(len(l.ip4)) {
			shuffleProxy[i] = l.ip4[j]
		}

	case Ip6:
		l.ip6Mutex.Lock()
		defer l.ip6Mutex.Unlock()

		// случайное прокси
		rand.Seed(time.Now().UnixNano())
		shuffleProxy = make([]*proxy, len(l.ip6))
		for i, j := range rand.Perm(len(l.ip6)) {
			shuffleProxy[i] = l.ip6[j]
		}

	default:
		return nil
	}

	for _, item := range shuffleProxy {
		if item.isFree(resource) {
			item.setBusy(resource)
			return item.url
		}
	}

	return nil
}

// GetFreeIP4 returns a free proxy Ip4 and marks it as busy.
func (l *List) GetFreeIP4(resource string) *url.URL {
	return l.GetFree(resource, Ip4)
}

// GetFreeIP6 returns a free proxy Ip6 and marks it as busy.
func (l *List) GetFreeIP6(resource string) *url.URL {
	return l.GetFree(resource, Ip6)
}

// SetFree removes the busy flag from the proxy. typeIP - proxy type (Ip4 or Ip6).
func (l *List) SetFree(resource string, u *url.URL, typeIP int) {
	var index int
	if index = l.Index(u, typeIP); index < 0 {
		return
	}

	switch typeIP {
	case Ip4:
		l.ip4Mutex.Lock()
		l.ip4[index].setFree(resource)
		l.ip4Mutex.Unlock()

	case Ip6:
		l.ip6Mutex.Lock()
		l.ip6[index].setFree(resource)
		l.ip6Mutex.Unlock()
	}
}

// SetFreeIP4 removes the busy flag from the proxy IP4.
func (l *List) SetFreeIP4(resource string, u *url.URL) {
	l.SetFree(resource, u, Ip4)
}

// SetFreeIP6 removes the busy flag from the proxy IP6.
func (l *List) SetFreeIP6(resource string, u *url.URL) {
	l.SetFree(resource, u, Ip6)
}

// NumIP4 returns the total number of proxies.
func (l *List) NumIP4() int {
	l.ip4Mutex.RLock()
	defer l.ip4Mutex.RUnlock()
	return len(l.ip4)
}

// NumIP6 returns the total number of proxies.
func (l *List) NumIP6() int {
	l.ip6Mutex.RLock()
	defer l.ip6Mutex.RUnlock()
	return len(l.ip6)
}

// NumBusyIP4 returns the number of busy proxies IP4.
func (l *List) NumBusyIP4(resource string) (numBusy int) {
	l.ip4Mutex.RLock()
	for _, item := range l.ip4 {
		if item.isBusy(resource) {
			numBusy++
		}
	}
	l.ip4Mutex.RUnlock()
	return
}

// NumBusyIP6 returns the number of busy proxies IP6.
func (l *List) NumBusyIP6(resource string) (numBusy int) {
	l.ip6Mutex.RLock()
	for _, item := range l.ip6 {
		if item.isBusy(resource) {
			numBusy++
		}
	}
	l.ip6Mutex.RUnlock()
	return
}

// NumFreeIP4 returns the number of free proxies.
func (l *List) NumFreeIP4(resource string) (numFree int) {
	return l.NumIP4() - l.NumBusyIP4(resource)
}

// NumFreeIP6 returns the number of free proxies.
func (l *List) NumFreeIP6(resource string) (numFree int) {
	return l.NumIP6() - l.NumBusyIP6(resource)
}

// parsing parses proxy list.
func (l *List) parsing(b []byte) (good []*proxy, bad []string) {
	var (
		u   *url.URL
		err error
	)

	bad = []string{}
	good = []*proxy{}
	for _, row := range bytes.Split(b, []byte("\n")) {
		row = bytes.TrimSpace(row)
		if len(row) == 0 {
			continue
		}

		u, err = url.Parse(string(row))
		if err != nil || u.Host == "" {
			bad = append(bad, string(row))
			continue
		}

		good = append(good, newProxy(u))
	}
	return
}

// refresh updates the proxy list. typeIP - proxy type (Ip4 or Ip6).
func (l *List) refresh(b []byte, typeIP int) (bad []string, err error) {
	var good []*proxy

	if typeIP != Ip4 && typeIP != Ip6 {
		err = errors.New("Unknown type. ")
		return
	}

	if good, bad = l.parsing(b); len(good) == 0 {
		return
	}

	switch typeIP {
	case Ip4:
		l.refreshIP4(good)

	case Ip6:
		l.refreshIP6(good)
	}
	return
}

// refresh updates the proxy IP4 list.
func (l *List) refreshIP4(p []*proxy) {
	l.ip4Mutex.Lock()
	l.ip4 = p
	l.ip4Mutex.Unlock()
}

// refresh updates the proxy IP6 list.
func (l *List) refreshIP6(p []*proxy) {
	l.ip6Mutex.Lock()
	l.ip6 = p
	l.ip6Mutex.Unlock()
}

// isBusy returns whether the proxy is busy.
// Always returns false if no proxy is found or unknown type.
func (l *List) isBusy(resource string, u *url.URL, typeIP int) bool {
	var index int
	// l.Index -1 if unknown type.
	if index = l.Index(u, typeIP); index < 0 {
		return false
	}

	if typeIP == Ip4 {
		l.ip4Mutex.RLock()
		defer l.ip4Mutex.RUnlock()
		return l.ip4[index].isBusy(resource)
	}

	// Ip6
	l.ip6Mutex.RLock()
	defer l.ip6Mutex.RUnlock()
	return l.ip6[index].isBusy(resource)
}

// isBusyIP4 returns whether the proxy is busy.
// Always returns false if no proxy is found or unknown type.
func (l *List) isBusyIP4(resource string, u *url.URL) bool {
	return l.isBusy(resource, u, Ip4)
}

// isBusyIP6 returns whether the proxy is busy.
// Always returns false if no proxy is found or unknown type.
func (l *List) isBusyIP6(resource string, u *url.URL) bool {
	return l.isBusy(resource, u, Ip6)
}

// setBusy sets the busy flag to the proxy.
func (l *List) setBusy(resource string, u *url.URL, typeIP int) {
	var index int
	if index = l.Index(u, typeIP); index < 0 {
		return
	}

	switch typeIP {
	case Ip4:
		l.ip4Mutex.RLock()
		defer l.ip4Mutex.RUnlock()
		l.ip4[index].setBusy(resource)

	case Ip6:
		l.ip6Mutex.RLock()
		defer l.ip6Mutex.RUnlock()
		l.ip6[index].setBusy(resource)
	}

	return
}

// setBusyIP4 sets the busy flag to the proxy IP4.
func (l *List) setBusyIP4(resource string, u *url.URL) {
	l.setBusy(resource, u, Ip4)
}

// setBusyIP6 sets the busy flag to the proxy IP6.
func (l *List) setBusyIP6(resource string, u *url.URL) {
	l.setBusy(resource, u, Ip6)
}

// Index returns the proxy index in the array.
// returns -1 if not found.
func (l *List) Index(u *url.URL, typeIP int) (index int) {
	switch typeIP {
	case Ip4:
		l.ip4Mutex.RLock()
		defer l.ip4Mutex.RUnlock()

		for i, p := range l.ip4 {
			if p.Equal(u) {
				return i
			}
		}

	case Ip6:
		l.ip6Mutex.RLock()
		defer l.ip6Mutex.RUnlock()

		for i, p := range l.ip6 {
			if p.Equal(u) {
				return i
			}
		}
	}

	return -1
}

func (l *List) IndexIP4(u *url.URL) (index int) {
	return l.Index(u, Ip4)
}

func (l *List) IndexIP6(u *url.URL) (index int) {
	return l.Index(u, Ip6)
}
