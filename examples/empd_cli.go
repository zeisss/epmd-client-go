package main

import "github.com/zeisss/epmd-client-go" epmd
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
  log.Println("NodeName\t Port")
  for _, v := range names {
    log.Println(v.Name+"\t", v.Port)
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
