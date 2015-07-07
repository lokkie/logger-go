package main
import (
    "net"
    "time"
    "io"
    "fmt"
    "bufio"
    "errors"
    "encoding/json"
)

type (


    SlaveClient struct {
        serverAddr string
        conn net.Conn
    }
)
func NewSlave(serverAddr string) *SlaveClient {
    client := &SlaveClient{serverAddr:serverAddr}
    return client
}

func (client *SlaveClient) pushExisting(hash string, timestamp int) error {
    var err error
    var message []byte
    var response string
    message, err = json.Marshal(&ExistingRecord{hash:hash, timestamp:timestamp})
    if err != nil {
        return err
    }
    response, err = client.say("INC", string(message))
    if response == "OK" {
        return nil
    }
    return err
}

func (client *SlaveClient) pushNew(scope, tag, err, fileName, extended string, timestamp, lineIndex int) error {

    return nil
}

func (client *SlaveClient) say(command, data string) (response string, err error) {
    err = client.keepConnected()
    if err != nil {
        return "", err
    }
    // Write calculated delimiter algorithm
    delim:="\n"
    message := command+" "+delim+" "+data+delim
    fmt.Fprint(client.conn, message)
    response, err = bufio.NewReader(client.conn).ReadString('\n')
    if err != nil {
        return "", err
    }
    return response, err
}

func (client *SlaveClient) keepConnected() error {
    if !client.connected() {
        iteration := 1
        for {
            err := client.connect()
            if err != nil {
                if (iteration > 60) {
                    return err
                }
                time.Sleep(1000 * time.Millisecond * iteration)
                iteration++
            }
        }
    }
    return nil
}

func (client *SlaveClient) connected() bool {
    if client.conn == nil {
        return false
    }
    client.conn.SetReadDeadline(time.Now())
    var one []byte
    _, err := client.conn.Read(one)
    if  err == io.EOF {
        client.cleanUp()
        return false
    } else {
        client.conn.SetReadDeadline(time.Time{})
    }
    if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
        client.cleanUp()
        return false
    }
    return true
}

func (client *SlaveClient) cleanUp() {
    if client.conn != nil {
        client.conn.Close();
        client.conn = nil
    }
}

func (client *SlaveClient) connect() error {
    var err error
    client.conn, err = net.Dial("tcp", client.serverAddr)
    if err != nil {
        return err
    }
    defer client.cleanUp()
    return client.handshake()
}

func (client *SlaveClient) handshake() error {
    fmt.Fprintf(client.conn, "MODE SYNC"+MESSAGE_DELIMITER)
    message, err := bufio.NewReader(client.conn).ReadString(MESSAGE_DELIMITER)
    if err != nil {
        return err
    }
    if message != "OK" {
        return errors.New("Could't do handshake")
    }
    return nil
}