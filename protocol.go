package main
import (
    "encoding/json"
    "io"
    "bufio"
    "regexp"
    "errors"
    "strings"
    "fmt"
)

type (
    Mode int

    NewRecord struct {
        scope string `json:"scope"`
        tag string `json:"tag"`
        err string `json:"exception"`
        errExtended string `json:"exception-text"`
        fileName string `json:"file"`
        lineIndex int `json:"line"`
        timestamp int `json:"timestamp"`
        extended string `json:"extended"`
    }

    ExistingRecord struct {
        hash string `json:"item-hash"`
        timestamp int `json:"timestamp"`
    }

    Message struct {
        raw, command string
        data interface{}
    }
    
    LoggerProtocol struct {
    }
)

const {
    SYNC Mode = iota
    ASYNC 
}

const (
    MESSAGE_DELIMITER = "\x00"
)

func MessageFromString(str string) (*Message, error) {
    message := &Message{raw:str}
    return message, message.unpack()
}

func NewMessage(command string, data interface{}) (*Message, error) {
    message := &Message{command:command, data:data}
    err := message.pack()
    if err != nil {
        return nil, err
    }
    return message, nil
}

func MessageFromStream(rd io.Reader) (*Message, error) {
    raw, err := bufio.NewReader(rd).ReadString(MESSAGE_DELIMITER)
    if err != nil {
        return nil, err
    }
    return MessageFromString(raw)
}



func (message *Message) unpack() error{
    dataCommandRegEx, err := regexp.Compile("(.*?) (.*)"+MESSAGE_DELIMITER)
    if err != nil {
        return err
    }
    commandOnlyRegEx, err := regexp.Compile("(.*)\x00")
    if err != nil {
        return err
    }
    if dataCommandRegEx.MatchString(message.raw) {
        parts := dataCommandRegEx.FindAllStringSubmatch(message.raw, -1)[0][1:]
        message.command = parts[0]
        message.data = CreateDataInterface(message.command, parts[1])
    } else if commandOnlyRegEx.MatchString(message.raw) {
        message.command = commandOnlyRegEx.FindString(message.raw)
        message.data = nil
    } else {
        return errors.New("No valid command found")
    }
    return nil
}


// Message structure:
// COMMAND DATA\0
func (message *Message) pack() error {
    data, err := json.Marshal(message.data)
    if err != nil {
        return err
    }
    message.raw = message.command + " " + string(data) + MESSAGE_DELIMITER
    return nil
}

func CreateDataInterface(command, data string) interface{} {
    var iData interface{}
    switch strings.ToUpper(command) {
        default:
            iData = nil
        case "NEW":
            iData = &NewRecord{}
        case "APPEND":
            iData = &ExistingRecord{}
    }
    if iData != nil {
        return json.Unmarshal([]byte(data), iData)
    } else {
        return nil
    }
}



func (lp *LoggerProtocol) doClientHandshake(conn net.Conn, mode Mode) (bool, error) {
    var message string = "MODE " 
    if mode = SYNC {
        message += "SYNC"+MESSAGE_DELIMITER
    } else if mode == ASYNC {
        message += "ASYNC"+MESSAGE_DELIMITER
    } else {
        return false, errors.New("Unknown mode")
    }
    fmt.Fprintf(conn, message)
    message, err := bufio.NewReader(conn).ReadString(MESSAGE_DELIMITER)
    if err != nil {
        return err
    }
    if message != "OK" {
        return errors.New("Could't do handshake")
    }
    return nil
}