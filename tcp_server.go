package main

import (
    "container/list"
    "net"
    "strings"
)

func (c *Client) Read(buffer []byte) bool {
    _, error := c.Conn.Read(buffer)
    if error != nil {
        c.Close()
        return false
    }
    return true
}

func (c *Client) Close() {
    c.Quit <- true
    c.Conn.Close()
    c.RemoveMe()
}

func (c *Client) Equal(other *Client) bool {
    if c.Conn == other.Conn {
        return true
    }
    return false
}

func (c *Client) RemoveMe() {
    for entry := c.ClientList.Front(); entry != nil; entry = entry.Next() {
        client := entry.Value.(Client)
        if c.Equal(&client) {
            Log("RemoveMe: ", c.Conn.RemoteAddr())
            c.ClientList.Remove(entry)
        }
    }
}

func IOHandler(Incoming <-chan string, clientList *list.List) {
    for {
        input := <-Incoming
        sinput := strings.Split(input,";")
        for e := clientList.Front(); e != nil; e = e.Next() {
            client := e.Value.(Client)
            for _,ctarget := range client.Targets {
                if ctarget == sinput[0] {
                    client.Incoming <-input
                }
            }
        }
    }
}

func ClientReader(client *Client) {
    buffer := make([]byte, 2048)
    for client.Read(buffer) {
        send := string(buffer)
        client.Incoming <- send
        for i := 0; i < 2048; i++ {
            buffer[i] = 0x00
        }
    }
}

func ClientSender(client *Client) {
loop:
    for {
        select {
            case buffer := <-client.Incoming:
                send:=ConcatStr("",buffer,"\n")
                client.Conn.Write([]byte(send))
            case <-client.Quit:
                client.Conn.Close()
                break loop
            }
    }
}

func ClientHandler(conn net.Conn, ch chan string, clientList *list.List, targets map[string]TargetT) {
    buffer := make([]byte, 2048)
    bytesRead, error := conn.Read(buffer)
    if error != nil {
        Log("Client connection error: ", error)
    }
    income := string(buffer[0:bytesRead])
    ctargets := strings.Split(income, " ")
    newClient := &Client{ctargets, make(chan string), ch, conn, make(chan bool), clientList}
    go ClientSender(newClient)
    go ClientReader(newClient)
    clientList.PushBack(*newClient)
    for _,target := range ctargets {
        tmp:=targets[target]
        if tmp.List!=nil {
            newClient.Incoming <-tmp.GetMlist(target)
        }
        newClient.Incoming <-tmp.GetTarget(target)
    }
    Log("ClientHandler: ", newClient.Conn.RemoteAddr())
}

func tcpServer(in chan string, targets map[string]TargetT, config *ConfigT) {
    clientList := list.New()
    go IOHandler(in, clientList)
    tcpAddr, error := net.ResolveTCPAddr("tcp", config.Main.TcpBind)
    if error != nil {
        Log("Error: Could not resolve address") 
    } else {
        netListen, error := net.Listen(tcpAddr.Network(), tcpAddr.String()) 
        if error != nil {
            Log("Listner: " ,error)
        } else {
            defer netListen.Close()
            for {
                connection, error := netListen.Accept()
                if error != nil {
                    Log("Client error: ", error)
                } else {
                    go ClientHandler(connection, in, clientList, targets)
                }
            }
        }
    }
}
