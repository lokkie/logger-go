package main
import (
    "flag"
    "path/filepath"
    "os"
    "net"
    "log"
    "time"
    "os/signal"
    "syscall"
)

func main() {
    // Command-line configuration
    mode := *flag.String("mode", "master", "Run mode: master/slave.")
    connType := *flag.String("conn-type", "unix", "Connection type: unix for unix socket; net for tcp/ip")
    confPath := *flag.String("config-path", "%app%/conf/main.json", "Path to config file. %app% - application dir")
    flag.Parse()
    say("Starting in mode:", mode)

    // Loading config
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    pHolder := &Placeholder{}
    pHolder.add("%app%", dir)
    config := &Configuration{}
    e := config.Load(pHolder.parse(confPath))
    if e != nil {
        say(err)
        return
    }

    StartServices(connType, mode, config)

    for {
        time.Sleep(200)
    }
}

func ConnectionHolder(c net.Conn) {

}

func say(a ...interface{}) {
    log.Println(a...)
}

func SlaveRoutine(config *Configuration) {
    client := NewSlave(config.Communication.MasterAddress)
    for {

    }
}

func StartServices(connType, mode string, config *Configuration) {
    // Run master-client communication
    if mode == "slave" {
        go SlaveRoutine(config)
    }

    // Run logger communication
    var unixServ, tcpServ *Server
    if connType == "unix" {
        say("Starting socket")
        unixServ = NewSocketServer(config.Communication.UnixSocketPath, ConnectionHolder)
        go unixServ.Listen()
    }

    if connType == "tcp" || mode == "master" {
        say("Starting tcp")
        tcpServ = NewTCPServer(config.Communication.ServerListenAdd, config.Communication.ServerListenPort, ConnectionHolder)
        go tcpServ.Listen()
    }

    // Teminate routine
    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    go func() {
//        s := <-sigc
        say("Stopping service")
        if unixServ != nil {
            unixServ.stopListening()
        }
        if tcpServ != nil {
            tcpServ.stopListening()
        }
        os.Exit(0)
    }()
}
