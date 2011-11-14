package main

import (
    "web"
    "strings"
    "godis"
    "fmt"
    "os"
)

const(
    // characters used for short-urls
    SYMBOLS = "0123456789abcdefghijklmnopqrsuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ"
    // special key in redis, that is our global counter
    COUNTER = "__counter__"
    HTTP = "http"
)

// connecting to redis on localhost, db with id 0 and no password
var (
    redis *godis.Client
)

func bootstrap(cfg *Config) os.Error{
    host := cfg.GetStringDefault("redis.address", "tcp:localhost:6379")
    db := cfg.GetIntDefault("redis.database", 0)
    passwd := cfg.GetStringDefault("redis.password", "")

    println(host)
    println(db)
    println(passwd)

    redis = godis.New(host, db, passwd)
    return nil
}


// function to resolve a shorturl and redirect
func resolve(ctx *web.Context, short string) {
    redirect, err := redis.Get(short)
    if err == nil {
        ctx.Redirect(301, redirect.String())
    } else {
        ctx.Redirect(301, "https://www.youtube.com/watch?v=jRHmvy5eaG4")
    }
}

// function to shorten and store a url
func shorten(ctx *web.Context, data string){
    const jsntmpl = "{\"url\" : \"%s\", \"longurl\" : \"%s\"}\n"
    if url, ok := ctx.Request.Params["url"]; ok{
        if ! strings.HasPrefix(url, HTTP){
            url = fmt.Sprintf("%s://%s", HTTP, url)
        }
        ctr, _ := redis.Incr(COUNTER)
        encoded := encode(ctr)
        go redis.Set(encoded, url)
        request := ctx.Request
        ctx.SetHeader("Content-Type", "application/json", true)
        host := request.Host
        if realhost, ok := ctx.Request.Params["X-Real-IP"]; ok{
            host = realhost
        }
        location := fmt.Sprintf("%s://%s/%s", HTTP, host, encoded)
        ctx.SetHeader("Location", location, true)
        ctx.StartResponse(201)
        ctx.WriteString(fmt.Sprintf(jsntmpl, location, url))
    }else{
       ctx.Redirect(404, "/")
    }
}

// encodes a number into our *base* representation
// TODO can this be made better with some bitshifting?
func encode(number int64) string{
    const base = int64(len(SYMBOLS))
    rest := number % base
    // strings are a bit weird in go...
    result := string(SYMBOLS[rest])
    if number - rest != 0{
       newnumber := (number - rest ) / base
       result = encode(newnumber) + result
    }
    return result
}

// main function that inits the routes in web.go
func main() {
    cfg := NewConfig("conf/kurz.conf")
    cfg.Parse()
    err := bootstrap(cfg)
    if err == nil {
        // this could go to bootstrap as well
        web.Post("/shorten/(.*)", shorten)
        web.Get("/(.*)", resolve)
        listen := cfg.GetStringDefault("listen", "0.0.0.0")
        port := cfg.GetStringDefault("port", "9999")
        web.Run(fmt.Sprintf("%s:%s", listen, port))
    }
}

