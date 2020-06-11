package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const (
	// https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
	defDSN = "user=matrix_synapse password=PASS dbname=matrix_synapse"
)

func main() {
	ctx := context.Background()

	dsn := flag.String("dsn", defDSN,
		"postgres database connection string")
	roomID := flag.String("room", "",
		"matrix unencrypted room ID to scan messages")

	flag.Parse()

	if *roomID == "" {
		log.Fatal("-room option needs to be set")
	}

	if err := roomMsg(ctx, *dsn, *roomID); err != nil {
		log.Fatal(err)
	}
}

type MsgType string

const (
	Text  = MsgType("m.text")
	Image = MsgType("m.image")
	Video = MsgType("m.video")
	Audio = MsgType("m.audio")
)

type Msg struct {
	Content MsgContent `json:"content"`
}

func (m Msg) String() string {
	return m.Content.String()
}

type MsgContent struct {
	Typ  MsgType `json:"msgtype"`
	Body string  `json:"body"`
	URL  string  `json:"url"`
}

func (mc MsgContent) String() string {
	switch mc.Typ {
	case Text:
		return mc.Body
	case Image:
		return mc.Body
	case Video:
		return mc.Body
	case Audio:
		return mc.Body
	default:
		return "UNKNOWN Message Type"
	}
}

func roomMsg(ctx context.Context, dsn, roomID string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	defer db.Close()

	rows, err := db.QueryContext(ctx, `
SELECT events.origin_server_ts, events.sender, event_json.json
FROM events INNER JOIN event_json ON events.event_id = event_json.event_id
WHERE events.type = 'm.room.message' AND events.room_id = $1
ORDER BY events.stream_ordering;
`, roomID)

	if err != nil {
		return fmt.Errorf("error querying database: %w", err)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error querying database: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			ts     int64
			sender string
			jsn    string
			msg    Msg
		)

		if err := rows.Scan(&ts, &sender, &jsn); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}

		// Convert millisecond unix time to Go time.Time
		t := time.Unix(ts/1000, ts%1000*1000000)

		if err := json.Unmarshal([]byte(jsn), &msg); err != nil {
			return fmt.Errorf("error unmarshalling json: %w", err)
		}

		fmt.Printf("### %s - %s\n\n%s\n\n",
			sender,
			t.Format(time.RFC3339),
			msg,
		)
	}

	return nil
}
