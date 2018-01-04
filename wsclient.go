package dgo2poc

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"golang.org/x/oauth2"
)

var (
	// Returned if you call Open() on an already connected WSClient.
	ErrWSAlreadyOpen = errors.New("websocket connection is already open")

	// Returned if discord sends a WSOPInvalidSession.
	ErrWSInvalidSession = errors.New("session is invalid, try again later")
)

// WSClient is a client for the Discord websocket API.
type WSClient interface {
	// Maintains a connection to the WSAPI until the context is cancelled.
	Run(ctx context.Context) error

	// Send an arbitrary packet.
	Send(op WSOP, data interface{}) error

	// Adds event handler(s). Handlers are created by each event's On... function.
	// For example, to handle Ready events, use OnReady.
	AddHandler(hl ...wsHandler) func()
}

type wsClient struct {
	REST  Client
	Token *oauth2.Token
	Opts  []WSOpt

	SessionID string       // last session id, for resume
	Seq       atomic.Int64 // last seq received
	Handlers  wsHandlers   // all registered handlers

	send chan<- wsPayload // use with Send() wrapper
	recv chan wsPayload   // only access for testing!!
}

func NewWSClient(cl Client, opts ...WSOpt) WSClient {
	return &wsClient{REST: cl, Token: cl.Token(), Opts: opts}
}

func (c *wsClient) AddHandler(hls ...wsHandler) func() {
	switch len(hls) {
	case 0:
		return func() {}
	case 1:
		return hls[0](&c.Handlers)
	default:
		removers := make([]func(), len(hls))
		for i, hl := range hls {
			removers[i] = hl(&c.Handlers)
		}
		return func() {
			for _, fn := range removers {
				fn()
			}
		}
	}
}

func (c *wsClient) Run(ctx context.Context) (rerr error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gw, err := c.REST.Gateway(ctx)
	if err != nil {
		return err
	}
	gwURL := gw.URL + "?encoding=json&v=" + GatewayVersion

	wsDialer := websocket.Dialer{
		NetDial: func(network, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, addr)
		},
	}
	conn, res, err := wsDialer.Dial(gwURL, nil)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()
	log.Printf("wsclient: connected (%d)!", res.StatusCode)

	errC := make(chan error)
	recv := make(chan wsPayload)
	send := make(chan wsPayload)

	go func() { errC <- errors.Wrap(wsRecv(ctx, conn, recv), "recv") }()
	go func() { errC <- errors.Wrap(wsSend(ctx, conn, send), "send") }()
	go func() { errC <- errors.Wrap(c.run(ctx, recv, send), "run") }()

	for i := 0; i < 3; i++ {
		err := <-errC
		if err != nil && rerr == nil {
			rerr = err
		}
		cancel()
	}

	return
}

func (c *wsClient) run(ctx context.Context, recv, send chan wsPayload) error {
	c.recv = recv
	c.send = send
	defer func() {
		c.recv = nil
		c.send = nil
	}()

	// Register handlers, deregister them in a defer.
	defer c.AddHandler(
		OnReady(func(ctx context.Context, ev *Ready) {
			c.SessionID = ev.SessionID
		}),
	)()

	var heartbeat *time.Ticker
	for {
		var beat <-chan time.Time
		if heartbeat != nil {
			beat = heartbeat.C
		}

		select {
		case <-beat:
			log.Printf("wsclient: sending heartbeat...")
			if err := c.sendHeartbeat(); err != nil {
				return err
			}
		case pl := <-recv:
			switch pl.OP {
			case WSOPDispatch:
				log.Printf("wsclient: dispatch: %s: %s", pl.Type, string(pl.Data))
				if err := c.Handlers.Dispatch(ctx, pl.Type, pl.Data); err != nil {
					return errors.Wrapf(err, "%s", pl.Type)
				}
			case WSOPHello:
				log.Printf("wsclient: received hello...")
				// This will be received once, right after connecting to the gateway.
				// This could be generalised into an event, not sure if it'd make sense.
				var d struct {
					HeartbeatInterval int64 `json:"heartbeat_interval"`
				}
				if err := json.Unmarshal(pl.Data, &d); err != nil {
					return err
				}

				// Start the heartbeat timer.
				beat := time.Duration(d.HeartbeatInterval) * time.Millisecond
				heartbeat = time.NewTicker(beat)
				defer heartbeat.Stop()
				log.Printf("wsclient: heartbeat interval: %v\n", beat)

				// Respond with a WSOPIdentify.
				log.Printf("wsclient: identifying...")
				if err := c.sendIdentify(); err != nil {
					return err
				}
			case WSOPHeartbeat:
				// The server may send a WSOPHeartbeat to immediately request a beat.
				log.Printf("wsclient: heartbeat requested")
				if err := c.sendHeartbeat(); err != nil {
					return err
				}
			case WSOPHeartbeatAck:
				// TODO: If no ACK is received between heartbeats, reset + resume the connection!
			case WSOPInvalidSession:
				log.Printf("wsclient: invalid session")
				var resumable bool
				if err := json.Unmarshal(pl.Data, &resumable); err != nil {
					return err
				}
				if resumable {
					// TODO: Try to reset the connection!
				}
				return ErrWSInvalidSession
			case WSOPReconnect:
				log.Printf("wsclient: reconnect")
				// TODO: Reset the connection!
			default:
				log.Printf("unknown OP: %d (t=%s, s=%d, d=%s)", pl.OP, pl.Type, pl.Seq, string(pl.Data))
				return errors.Errorf("unknown OP: %d (t=%s, s=%d, d=%s)", pl.OP, pl.Type, pl.Seq, string(pl.Data))
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *wsClient) Send(op WSOP, d interface{}) error {
	pl := wsPayload{OP: op}
	if d != nil {
		data, err := json.Marshal(d)
		if err != nil {
			return err
		}
		pl.Data = json.RawMessage(data)
	}
	c.send <- pl
	return nil
}

func (c *wsClient) sendHeartbeat() error {
	var d interface{}
	if seq := c.Seq.Load(); seq != 0 {
		d = seq
	}
	return c.Send(WSOPHeartbeat, d)
}

func (c *wsClient) sendIdentify() error {
	opts := WSOpts{
		id: wsIdentify{
			Token:          c.Token.AccessToken,
			Compress:       false, // not yet supported :(
			LargeThreshold: 50,
			Shard:          [2]int{0, 1},
			Presence:       WSStatus{Status: WSStatusOnline},
			Properties: wsIdentifyProps{
				OS:      runtime.GOOS,
				Browser: "dgo2poc",
				Device:  "dgo2poc",
			},
		},
	}
	for _, opt := range c.Opts {
		opt(&opts)
	}
	return c.Send(WSOPIdentify, opts.id)
}
