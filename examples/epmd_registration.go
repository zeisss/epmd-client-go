package main

import epmd ".." // "github.com/zeisss/epmd-client-go"
import "log"

func main() {
  conn, err := epmd.Register(12345, epmd.NODE_TYPE_ERLANG, 5, 5, "demo", "")
  if err != nil {
    log.Fatalf("Failed to register with EPMD: %v", err)
  }
  defer conn.Close()

  log.Println("Process registered with EPMD. Use other tools to verify.")
  log.Println("Press ^C to close")
  select {}
}
