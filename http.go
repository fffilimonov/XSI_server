package main

import (
    "bytes"
    "io/ioutil"
    "net"
    "net/http"
    "time"
)

func handlePost (bodychan chan []byte) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        body,_ := ioutil.ReadAll(req.Body)
        bodychan<-body
    }
}

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, time.Duration(time.Second*2))
}

func sendPost(Event string, config *ConfigT) error {
    url := ConcatStr("","http://",config.Main.RemoteHost,"/com.broadsoft.xsi-events/v2.0/system")
    var xmlStr string = ConcatStr("",
                                "<?xml version=\"1.0\" encoding=\"utf-8\"?><Subscription xmlns=\"http://schema.broadsoft.com/xsi\"><event>",
                                Event,
                                "</event><httpContact><uri>",
                                config.Main.HttpContact,
                                "/events/system</uri></httpContact><applicationId>",
                                ConcatStr("",config.Main.AppID,Event),
                                "</applicationId></Subscription>")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(xmlStr)))
    if err == nil {
        req.SetBasicAuth(config.Main.User, config.Main.Password)
        transport := http.Transport{
            Dial: dialTimeout,
        }
        client := http.Client{
            Transport: &transport,
        }
        resp, err := client.Do(req)
        if err != nil {
            Log(err)
            return err
        } else {
            resp.Body.Close()
            Log("Subscribe Status:", resp.Status)
        }
    }
    return err
}

func httpServer (bodychan chan []byte, config *ConfigT) {
    var flag bool
    go func () {
        http.HandleFunc("/events/system", handlePost(bodychan))
        http.ListenAndServe(config.Main.HttpBind, nil)
    }()
    for {
        flag=false
        for _,event := range config.Main.Event {
            err:=sendPost(event, config)
            if err != nil {
                flag=true
            }
        }
        if flag {
            time.Sleep(time.Second*1)
        } else {
            time.Sleep(time.Second*time.Duration(config.Main.Expires))
        }
    }
}

func httpClient (bodychan chan []byte, URL string) {
    for {
        select {
            case body := <-bodychan:
            req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
            if err == nil {
                client := &http.Client{}
                resp, err := client.Do(req)
                if err != nil {
                    Log(err)
                } else {
                    resp.Body.Close()
                    Log("Send Status:", resp.Status)
                }
            }
        }
    }
}
