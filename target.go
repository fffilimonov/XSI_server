package main

import "strconv"

func (tmp *TargetT) AddMCall(time string, addr string) {
    var tmpcall lCalls
    tmpcall.Time=time
    tmpcall.Addr=addr
    tmp.List.Value = tmpcall
    tmp.List = tmp.List.Next()
}

func (tmp *TargetT) GetMlist(target string) string {
    var i uint
    outcalls := ConcatStr(";",target,"calls")
    if tmp.List!=nil {
        for i=0;i<uint(tmp.List.Len());i++ {
            tmp.List = tmp.List.Prev()
            if tmp.List.Value != nil {
                tmpcall:=tmp.List.Value.(lCalls)
                outcalls = ConcatStr(";",outcalls,tmpcall.Time,tmpcall.Addr)
            } else {
                outcalls = ConcatStr(";",outcalls,"","")
            }
        }
    }
    return outcalls
}

func (tmp *TargetT) GetTarget(target string) string {
    return ConcatStr(";",target,"state",tmp.Hook,tmp.CCstatus,tmp.Addr,tmp.CallID,strconv.Itoa(tmp.Qcalls))
}
