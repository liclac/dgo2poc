package dgo2poc

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func wsRecv(ctx context.Context, conn *websocket.Conn, recv chan<- wsPayload) error {
	for {
		var pl wsPayload
		if err := conn.ReadJSON(&pl); err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return err
			}
		}
		recv <- pl
	}
}

func wsSend(ctx context.Context, conn *websocket.Conn, send <-chan wsPayload) error {
	for {
		select {
		case pl := <-send:
			data, err := json.Marshal(pl)
			if err != nil {
				return err
			}
			log.Printf("wsclient: sending: %s", string(data))
			if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
				select {
				case <-ctx.Done():
					return nil
				default:
					return err
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}
