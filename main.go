package main

import (
    "fmt"
    "os"
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
    targets := make(map[string]TargetT)
    in := make(chan string)

    go httpServer (bodychan, &config)
    go tcpServer (in, targets, &config)
    go parser (bodychan, in, targets)
    for {
        Log("in map: ",len(targets))
        time.Sleep(time.Second*time.Duration(config.Main.Expires))
    }
}
