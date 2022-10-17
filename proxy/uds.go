package proxy

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

const srvPattern = "/srv/register"
const srvUrl = "http://unix" + srvPattern

func RegisterSrv(server Server) error {
	logx.Infof("register srv: %+v\n", server)
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketAddress)
			},
		},
	}

	srv := httpc.NewServiceWithClient("uds-api", client)
	_, err := srv.Do(context.Background(), http.MethodPost, srvUrl, server)
	if err != nil {
		return err
	}

	return nil
}

func ListenSrv(ch *ConfigHandler, filters ...Filter) error {
	if err := os.RemoveAll(socketAddress); err != nil {
		return err
	}

	l, err := net.Listen("unix", socketAddress)
	if err != nil {
		return err
	}
	defer l.Close()

	var rh = &reqHandler{
		p: NewProxy(ch, filters...),
	}

	http.HandleFunc(srvPattern, rh.handleRequest)
	return http.Serve(l, nil)
}

type reqHandler struct {
	p *Proxy
}

func (rh *reqHandler) handleRequest(w http.ResponseWriter, r *http.Request) {
	logx.Info("---received a srv register event")
	log := logx.WithContext(r.Context())
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req Server
	if err := json.Unmarshal(body, &req); err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	go rh.p.Run(req)

	w.WriteHeader(http.StatusOK)
}
