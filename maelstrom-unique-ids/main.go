package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type writeMessageBody struct {
    maelstrom.MessageBody
    Id string `json:"id"`
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
    now := time.Now()
    return fmt.Sprintf("%s_%d", nid, now.UnixMicro())
}
