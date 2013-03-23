package main

import epmd ".." // "github.com/zeisss/epmd-client-go"
import "log"

func main() {
	client, err := epmd.Local()
	if err != nil {
		log.Fatalf("Failed to connect to epmd: %v", err)
	}

	names, err := client.Names()
	if err != nil {
		log.Fatalf("Failed to get name list from epmd: %v", err)
	}
	log.Println("NodeName\t Port\t Extra")
	for _, v := range names {
		// Fetch details
		node, err := client.Get(v.Name)
		if err != nil {
		  log.Println(v.Name+"\t", v.Port, "\t", err)
		} else {
		  log.Println(v.Name+"\t", v.Port, "\t", node.Extra)
		}
	}

	log.Println("")
	log.Println("================")
	log.Println("")

	node, err := client.Get("rabbit")
	log.Println("node:", node, "err:", err)
	log.Println("-----")

	node, err = client.Get("riak")
	log.Println("node:", node, "err:", err)

	log.Println("-----")

	node, err = client.Get("doesnotexist")
	log.Println("node:", node, "err:", err)
}
