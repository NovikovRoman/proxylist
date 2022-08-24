package proxylist

import (
	"net/url"
)

// proxy structure.
type proxy struct {
	listBusy map[string]bool
	url      *url.URL
}

func newProxy(u *url.URL) *proxy {
	return &proxy{
		url:      u,
		listBusy: map[string]bool{},
	}
}

func (p *proxy) String() string {
	return p.url.String()
}

func (p *proxy) Equal(u *url.URL) bool {
	return p.url.String() == u.String()
}

func (p *proxy) setBusy(resource string) {
	p.listBusy[resource] = true
}

func (p *proxy) isBusy(resource string) bool {
	return p.listBusy[resource]
}

func (p *proxy) setFree(resource string) {
	p.listBusy[resource] = false
}

func (p *proxy) isFree(resource string) bool {
	return !p.listBusy[resource]
}
