package main

import (
    "fmt"
    "os"
)

func Log (v ...interface{}) {
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
    parser (bodychan, in, targets)
}
