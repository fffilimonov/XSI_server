package main

import (
    "gopkg.in/gcfg.v1"
    "os"
)

func ReadConfig(Configfile string) ConfigT {
    var Config ConfigT
    err := gcfg.ReadFileInto(&Config, Configfile)
    if err != nil {
        Log(err)
        os.Exit (1)
    }
    return Config
}
