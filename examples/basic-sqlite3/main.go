package main

import (
	"fmt"
	"net/http"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
    "github.com/fiatjaf/khatru/policies"
)

func main() {
	relay := khatru.NewRelay()

	db := sqlite3.SQLite3Backend{DatabaseURL: "/tmp/khatru-sqlite-tmp"}
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

    allowedEventKinds := []uint16{37515, 33811, 13811, 30100, 31001, 34235, 34236}
	relay.RejectEvent = append(relay.RejectEvent, policies.RestrictToSpecifiedKinds(true, allowedEventKinds[0]))

    // Custom policy
    //relay.RejectEvent = append(relay.RejectEvent,
    //    // We only accept events with kind 37515, 13811 so we put them in an array
    //    func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
    //        fmt.Printf("%T: %d \n", event.Kind, event.Kind)
    //        slices.Sort(allowedEventKinds)
    //        n, found := slices.BinarySearch(allowedEventKinds, uint16(event.Kind))
    //        fmt.Println(n, found)
    //        if found {
    //            return false, ""
    //        }
    //        return true, "This event kind not allowed on this relay"
    //    },
    //)

    // Output when there is HTTP request
    mux := relay.Router()
    // set up other http handlers
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("content-type", "text/html")
        fmt.Fprintf(w, `<html><head></head><body>`)
        fmt.Fprintf(w, `<div style="text-align: center;">`)
        fmt.Fprintf(w, `Connect your Nostr client to <code>wss://khatru.nostrver.se</code>`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `This relay only accepts events with kind:`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `- <code>37515</code> (places)`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `- <code>33811, 13811</code> (check-ins)`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `- <code>34235, 34236</code> (NIP-71)`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `<code>30100, 30101</code> (draft NIP-113 activity events)`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `<a href="https://github.com/nostrver-se/khatru" target="_blank">https://github.com/nostrver-se/khatru</a>`)
        fmt.Fprintf(w, `</div>`)
        fmt.Fprintf(w, `</body></html>`)
    })

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}
