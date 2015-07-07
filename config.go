package main
import (
    "encoding/json"
    "io/ioutil"
)

type Configuration struct {
    Communication struct {
        MasterAddress    string  `json:"masterAddress"`
        ServerListenAdd  string  `json:"serverListenAdd"`
        ServerListenPort string `json:"serverListenPort"`
        UnixSocketPath   string  `json:"unixSocketPath"`
    } `json:"communication"`
}

func (conf *Configuration) Load(path string) error {
    jsonSrc, e := ioutil.ReadFile(path)
    if e != nil {
        return e
    }
    json.Unmarshal(jsonSrc, &conf)
    return nil
}