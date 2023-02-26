package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type writeMessageBody struct {
    maelstrom.MessageBody
    Id string `json:"id"`
}

type counter int32
var c counter

func (c *counter) inc() int32 {
    return atomic.AddInt32((*int32)(c), 1)
}

func (c *counter) get() int32 {
    return atomic.LoadInt32((*int32)(c))
}

func main() {
    n := maelstrom.NewNode()
    
    n.Handle("generate", func(msg maelstrom.Message) error {
        var reqBody maelstrom.MessageBody
        if err := json.Unmarshal(msg.Body, &reqBody); err != nil {
            return err
        }

	    nid := n.ID()
        var resBody writeMessageBody
        resBody.Type = "generate_ok"
        resBody.Id = generateId(nid)

        return n.Reply(msg, resBody)
    })

    if err := n.Run(); err != nil {
        log.Fatal(err)
    }
}

func generateId(nid string) string {
    count := c.get()
    c.inc()
    now := time.Now()
    return fmt.Sprintf("%d-%s-%d", now.UnixMilli(), nid, count)
}
