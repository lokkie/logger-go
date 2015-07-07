package main
import (
    "net"
    "os"
)

type Server struct {
    connectionType string
    path string
    connHandler func(c net.Conn)
    listener net.Listener
}

func NewTCPServer(listenAddr string, listenPort string, connHandler func(c net.Conn)) *Server {
    serv := &Server{
        connectionType:"tcp",
        path:listenAddr+":" +string(listenPort),
        connHandler:connHandler}

    return serv
}

func NewSocketServer(path string, connHandler func(c net.Conn)) *Server {
    srv := &Server{
        connectionType:"unix",
        path:path,
        connHandler:connHandler}
    return srv
}

func (serv *Server) Listen() error {
    var err error
    serv.listener, err = net.Listen(serv.connectionType, serv.path)
    if err != nil {
        say("Error listening:" + err.Error())
        return err
    }
    defer serv.stopListening()
    say("Started on: "+serv.path)
    for {
        conn, err := serv.listener.Accept()
        if err != nil {
            say("Error accepting" + err.Error())
            return err
        }
        go serv.connHandler(conn);
    }
}

func (serv *Server) stopListening() {
    say("Cleaning up: " + serv.connectionType)
    serv.listener.Close()
    if serv.connectionType == "unix" {
        os.Remove(serv.path)
    }
}