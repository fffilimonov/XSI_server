package main

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "time"
)

func handlePost (bodychan chan []byte) http.HandlerFunc {
    return func(w http.ResponseWriter, req *http.Request) {
        body,_ := ioutil.ReadAll(req.Body)
        bodychan<-body
    }
}

func sendPost(Event string, config *ConfigT) {
    url := ConcatStr("","http://",config.Main.RemoteHost,"/com.broadsoft.xsi-events/v2.0/system")
    var xmlStr string = ConcatStr("",
                                "<?xml version=\"1.0\" encoding=\"utf-8\"?><Subscription xmlns=\"http://schema.broadsoft.com/xsi\"><event>",
                                Event,
                                "</event><httpContact><uri>",
                                config.Main.HTTPServer,
                                "/events/system</uri></httpContact><applicationId>",
                                config.Main.AppID,
                                "</applicationId></Subscription>")
    req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(xmlStr)))
    req.SetBasicAuth(config.Main.User, config.Main.Password)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        Log(err)
    }
    defer resp.Body.Close()
    Log("response Status:", resp.Status)
}

func httpServer (bodychan chan []byte, config *ConfigT) {
    go func () {
        http.HandleFunc("/events/system", handlePost(bodychan))
        http.ListenAndServe(config.Main.HTTPBind, nil)
    }()
    for {
        for _,event := range config.Main.Event {
            sendPost(event, config)
        }
        time.Sleep(time.Second*600)
    }
}
