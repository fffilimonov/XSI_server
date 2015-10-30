package main

import (
    "container/list"
    "container/ring"
    "net"
)

type Client struct {
    Targets []string
    Incoming chan string
    Outgoing chan string
    Conn net.Conn
    Quit chan bool
    ClientList *list.List
}

type TargetT struct {
    Hook string
    CCstatus string
    Addr string
    CallID string
    Qcalls int
    List *ring.Ring
}

type EventT struct {
    Target string `xml:"targetId"`
    AppID string `xml:"externalApplicationId"`
    Edata edata `xml:"eventData"`
}

type edata struct {
    Etype string `xml:"type,attr"`
    Hook string `xml:"hookStatus"`
    Pers string `xml:"call>personality"`
    CCstatus string `xml:"agentStateInfo>state"`
    State string `xml:"call>state"`
    Cause string `xml:"call>releaseCause>internalReleaseCause"`
    Addr string `xml:"call>remoteParty>address"`
    CallID string `xml:"call>callId"`
    Qcalls int `xml:"position"`
    CCaddr string `xml:"queueEntry>remoteParty>address"`
    CCtime string `xml:"queueEntry>removeTime"`
    Rtime string `xml:"call>releaseTime"`
    Atime string `xml:"call>answerTime"`
}

type lCalls struct {
    Addr string
    Time string
}

type ConfigT struct {
    Main struct {
        User string
        Password string
        RemoteHost string
        HttpBind string
        HttpContact string
        Expires int
        Event []string
        AppID string
        TcpBind string
    }
    Reloadable struct {
        ASURL string
        StatsURL string
    }
}
