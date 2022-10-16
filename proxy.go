package proxy

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
)

const defaultProxyAddress = "0.0.0.0:9999"

const (
	StateInvalid State = iota
	StateConnected
	StateDisconnect
)

const (
	_ Signal = iota
	SignalHealthCheck
)

type State int

type Proxy struct {
	address               string
	handler               *Handler
	service               Service
	state                 State
	lock                  *sync.Mutex
	healthCheckOnce       sync.Once
	healthCheckTicker     *time.Ticker
	healthCheckTickerOnce sync.Once
	conn                  *grpc.ClientConn
}

type Option func(p *Proxy)

type Signal int

type Handler struct {
	ch     chan Service
	done   chan struct{}
	signal chan Signal
}

func NewHandler() *Handler {
	return &Handler{
		ch: make(chan Service),
	}
}

func (h *Handler) SendSignal(signal Signal) {
	h.signal <- signal
}

func WithProxyAddress(address string) Option {
	return func(p *Proxy) {
		p.address = address
	}
}

func NewProxy(handler *Handler, option ...Option) (*Proxy, error) {
	var p = &Proxy{
		handler:               handler,
		lock:                  &sync.Mutex{},
		healthCheckOnce:       sync.Once{},
		healthCheckTickerOnce: sync.Once{},
		healthCheckTicker:     time.NewTicker(500 * time.Millisecond),
	}

	for _, opt := range option {
		opt(p)
	}

	if len(p.address) == 0 {
		p.address = defaultProxyAddress
	}

	th := proxy.TransparentHandler(p.director)
	serviceHandlerOption := grpc.UnknownServiceHandler(th)
	s := grpc.NewServer(serviceHandlerOption)
	l, err := net.Listen("tcp", p.address)
	if err != nil {
		return nil, err
	}

	defer s.GracefulStop()
	if err = s.Serve(l); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Proxy) run() {
	go func() {
		for {
			select {
			case srv := <-p.handler.ch:
				p.updateService(srv)
			case <-p.handler.done:
				p.Close()
				return
			case sig := <-p.handler.signal:
				p.handleSignal(sig)
			}
		}
	}()
}

func (p *Proxy) handleSignal(signal Signal) {
	switch signal {
	case SignalHealthCheck:
		if err := p.service.Ping(); err != nil {
			p.updateState(StateDisconnect)
		} else {
			p.updateState(StateConnected)
		}
	}
}

func (p *Proxy) updateState(state State) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.state = state
}

func (p *Proxy) updateService(srv Service) {
	p.lock.Lock()
	p.service = Service{
		Name:    srv.Name,
		Address: srv.Address,
	}
	p.lock.Unlock()

	if err := srv.Ping(); err == nil {
		p.updateState(StateConnected)
	} else {
		p.updateState(StateDisconnect)
	}

	p.healthCheck()
}

func (p *Proxy) healthCheck() {
	p.healthCheckOnce.Do(func() {
		go func() {
			for {
				select {
				case <-p.healthCheckTicker.C:
					p.handler.SendSignal(SignalHealthCheck)
				case <-p.handler.done:
					return
				}
			}
		}()
	})
}

func (p *Proxy) Close() {
	if p.healthCheckTicker != nil {
		p.healthCheckTickerOnce.Do(func() {
			p.healthCheckTicker.Stop()
		})
	}
}

func (p *Proxy) director(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
	if p.conn != nil {
		return ctx, p.conn, nil
	}
	panic("implement me")
}
