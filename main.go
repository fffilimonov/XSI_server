package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func Log (v ...interface{}) {
    fmt.Fprintf(os.Stderr,"%s ",time.Now().Format(time.UnixDate))
    fmt.Fprint(os.Stderr,v...)
    fmt.Fprint(os.Stderr,"\n")
}

func main() {
    larg:=len(os.Args)
    if larg < 2 {
        Log("no args")
        os.Exit (1)
    }
    var confFile string = os.Args[1]
    config:=ReadConfig(confFile)

    bodychan := make(chan []byte)
    notifychan := make(chan []byte)
    targets := make(map[string]TargetT)
    in := make(chan string)

    go httpServer (bodychan, &config)
    go tcpServer (in, targets, &config)
    go parser (bodychan, in, targets, notifychan)
    go httpClient (notifychan, &config)

    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGHUP)
    for {
        select {
            case <-sig:
                Log("in map: ",len(targets))
                config=ReadConfig(confFile)
        }
    }
}
