package utils

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/proxy"
	"net"
	"strings"
)

var (
	ErrInvalidProxyFormat = errors.New("invalid socks5 format (must be [user:pass]@ip:port)")
)

type Proxy interface {
	Dialer() (proxy.Dialer, error)
	Conn(endpoint string) (net.Conn, error)
	Test(endpoint string) error
	String() string
}

type socks5 struct {
	addr string
	auth *proxy.Auth
}

func NewSocks5(raw string) (Proxy, error) {
	if len(raw) == 0 {
		return nil, ErrInvalidProxyFormat
	}
	p := &socks5{}

	args := strings.Split(raw, "@")
	if len(args) == 1 {
		p.addr = raw
	} else if len(args) > 1 {
		credentials := strings.Split(args[0], ":")
		if len(credentials) != 2 {
			return nil, ErrInvalidProxyFormat
		}
		p.auth = &proxy.Auth{
			User:     credentials[0],
			Password: credentials[1],
		}
		p.addr = args[1]
	}

	return p, nil
}

func (p *socks5) Dialer() (proxy.Dialer, error) {
	dialer, err := proxy.SOCKS5("tcp", p.addr, p.auth, proxy.Direct)
	if err != nil {
		return nil, errors.Wrap(err, "proxy.SOCKS5")
	}
	return dialer, nil
}

func (p *socks5) Conn(endpoint string) (net.Conn, error) {
	dialer, err := p.Dialer()
	if err != nil {
		return nil, errors.Wrap(err, "Dialer")
	}

	perHost := proxy.NewPerHost(dialer, proxy.Direct)
	if conn, err := perHost.Dial("tcp", endpoint); err != nil {
		return nil, errors.Wrap(err, "perHost.Dial")
	} else {
		return conn, nil
	}
}

func (p *socks5) Test(endpoint string) error {
	if conn, err := p.Conn(endpoint); err != nil {
		return err
	} else {
		conn.Close()
		return nil
	}
}

func (p *socks5) String() string {
	return fmt.Sprintf("addr=%s user=%s password=%s", p.addr, p.auth.User, p.auth.Password)
}
