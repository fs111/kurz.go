package main

import (
    "os"
    "strings"
    "io/ioutil"
    "strconv"
)


const (
    COMMENT = "#"
    NEWLINE = "\n" 
    SEPARATOR = "="
)


type Config struct{
    Path string
    Entries map[string]string
}


func NewConfig(path string) *Config{
    cfg := new(Config)
    cfg.Path = path
    cfg.Entries = make(map[string]string)
    return cfg
}

func (c Config) Parse() os.Error{
    contents, err := ioutil.ReadFile(c.Path)
    if err!= nil{
        return err
    }
    for _, line := range strings.Split(string(contents), "\n"){
        line = strings.TrimSpace(line)
        if int64(len(line)) > 2 && !strings.HasPrefix(line, COMMENT){
            row := strings.SplitN(line, SEPARATOR, 2)
            c.Entries[string(row[0])] = string(row[1])
        }
    }
    return nil
}

func (c Config) GetString(key string) (entry string, err os.Error){
    value, ok :=  c.Entries[key]
    var e os.Error
    if !ok{
        e = os.NewError("boom")
    } else {
        e = nil
    }
    return value, e
}



func (c Config) GetStringDefault(key string, val string) string{
    entry, e := c.GetString(key)
    if e == nil{
        return entry
    }
    return val
}


func (c Config) GetInt(key string) (val int, err os.Error){
    entry, e := c.GetString(key)
    if e == nil {
        return strconv.Atoi(entry)
    }
    // TODO figure out error handling
    return -1, os.NewError("eeep")
}

func (c Config) GetIntDefault(key string, val int) int{
    entry, e := c.GetInt(key)
    if e == nil {
        return entry
    }
    return val
}
