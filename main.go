package main

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/nbd-wtf/go-nostr"
)

var (
	relays = []string{
		"wss://nostr.band",
		"wss://yabu.me",
	}
)

func StreamHandler(w http.ResponseWriter, req *http.Request) {
	defer log.Println("closed")
	pool := nostr.NewSimplePool(context.TODO())

	flusher, _ := w.(http.Flusher)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan *nostr.Event)

	go func() {
		defer close(ch)

		now := nostr.Now()
		incomming := pool.SubMany(ctx, relays, nostr.Filters{
			{
				Kinds: []int{
					nostr.KindTextNote,
					//nostr.KindRepost,
					//nostr.KindReaction,
				},
				Since: &now,
			},
		})
		for {
			select {
			case <-req.Context().Done():
				return
			case ev := <-incomming:
				ch <- ev.Event
			case <-ctx.Done():
				return
			}
		}
	}()

	for ev := range ch {
		b, err := json.Marshal(ev)
		if err != nil {
			return
		}
		_, err = w.Write(b)
		if err != nil {
			return
		}
		log.Println("flush")
		flusher.Flush()
	}
}

var (
	//go:embed static
	assets embed.FS
)

func main() {
	sub, _ := fs.Sub(assets, "static")
	http.HandleFunc("/stream", StreamHandler)
	http.Handle("/", http.FileServer(http.FS(sub)))
	log.Println("listening :8080")
	http.ListenAndServe(":8080", nil)
}
