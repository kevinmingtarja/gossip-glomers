package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type requestMessageBody struct {
    maelstrom.MessageBody
    Message int `json:"message"`
}

type readResponseMessageBody struct {
    maelstrom.MessageBody
    Messages []int `json:"messages"`
}

type topologyRequestMessageBody struct {
    maelstrom.MessageBody
    Topology []int `json:"topology"`
}

var seen []int

func main() {
    n := maelstrom.NewNode()
    seen = []int{}
    
    n.Handle("broadcast", func(msg maelstrom.Message) error {
        var reqBody requestMessageBody
        if err := json.Unmarshal(msg.Body, &reqBody); err != nil {
            return err
        }

        storeMessage(reqBody.Message)

        nodeID := n.ID()
        nodeIDs := n.NodeIDs()
        for _, nid := range nodeIDs {
            if nid != nodeID {
                broadcast(nid, reqBody.Message)
            }
        }

        var resBody maelstrom.MessageBody
        resBody.Type = "broadcast_ok"

        return n.Reply(msg, resBody)
    })

    n.Handle("read", func(msg maelstrom.Message) error {
        var resBody readResponseMessageBody
        resBody.Type = "read_ok"
        resBody.Messages = seen

        return n.Reply(msg, resBody)
    })
    
    n.Handle("topology", func(msg maelstrom.Message) error {
        var reqBody requestMessageBody
        if err := json.Unmarshal(msg.Body, &reqBody); err != nil {
            return err
        }
        
        var resBody maelstrom.MessageBody
        resBody.Type = "topology_ok"

        return n.Reply(msg, resBody)
    })

    if err := n.Run(); err != nil {
        log.Fatal(err)
    }
}

func storeMessage(msg int) {
    seen = append(seen, msg)
}

func broadcast(to string, value int) {
    return
}
