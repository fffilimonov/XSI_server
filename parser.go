package main

import (
    "container/ring"
    "encoding/xml"
    "strings"
)

func ParseXml (data []byte) EventT {
    var edata EventT
    xml.Unmarshal(data, &edata)
    return edata
}

func ConcatStr (sep string, args ... string) string {
    return strings.Join(args, sep)
}

func parser (bodychan chan []byte, in chan string, targets map[string]TargetT, notifychan chan []byte, statschan chan []byte) {
    for {
        select {
            case body := <-bodychan:
                parsedBody := ParseXml(body)
                if parsedBody.Target!="" {
                    tmp:=targets[parsedBody.Target]

                    if parsedBody.AppID == "GolangClientCall Center Queue" ||
                     parsedBody.AppID == "GolangClientCall Center Agent" ||
                     tmp.CCstatus == "Available" {
                        statschan <- body
                    }

                    if tmp.List == nil {
                        tmp.List = ring.New(15)
                    }

                    if parsedBody.Edata.Hook!="" {
                        tmp.Hook=parsedBody.Edata.Hook
                    }

                    if parsedBody.Edata.CCstatus!="" {
                        tmp.CCstatus=parsedBody.Edata.CCstatus
                    }

                    if parsedBody.Edata.Pers == "Terminator" && parsedBody.Edata.State == "Alerting" {
                        tmp.Addr=parsedBody.Edata.Addr
                        tmp.CallID=parsedBody.Edata.CallID
                        notifychan <- body
                    }

                    if parsedBody.Edata.Pers == "Terminator" && parsedBody.Edata.State == "Released" {
                        tmp.Addr=""
                        tmp.CallID=""
                        if parsedBody.Edata.Atime == "" {
                            tmp.AddMCall(parsedBody.Edata.Rtime, parsedBody.Edata.Addr)
                            in <-tmp.GetMlist(parsedBody.Target)
                        }
                        notifychan <- body
                    }

                    if parsedBody.Edata.Pers == "Terminator" && parsedBody.Edata.State == "Active" {
                        notifychan <- body
                    }

                    if parsedBody.Edata.Etype == "xsi:ACDCallAddedEvent" {
                        tmp.Qcalls=parsedBody.Edata.Qcalls
                    }

                    if parsedBody.Edata.Etype == "xsi:ACDCallOfferedToAgentEvent" {
                        if tmp.Qcalls > 0 {
                            tmp.Qcalls--
                        }
                    }

                    if parsedBody.Edata.Etype == "xsi:ACDCallAbandonedEvent" {
                        if tmp.Qcalls > 0 {
                            tmp.Qcalls--
                        }
                        tmp.AddMCall(parsedBody.Edata.CCtime, parsedBody.Edata.CCaddr)
                        in <-tmp.GetMlist(parsedBody.Target)
                    }
                    in <-tmp.GetTarget(parsedBody.Target)
                    targets[parsedBody.Target]=tmp
                }
        }
    }
}
