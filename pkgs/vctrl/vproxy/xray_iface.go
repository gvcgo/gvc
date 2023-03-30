package vproxy

type Proxies interface {
	Today() string
	GetDate() string
	Reload()
	Update(any)
}
