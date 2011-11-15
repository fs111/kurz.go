package main

import (
    "web"
    "strings"
    "godis"
    "fmt"
    "os"
    "url"
    "flag"
)

const(
    // special key in redis, that is our global counter
    COUNTER = "__counter__"
    HTTP = "http"
)

// connecting to redis on localhost, db with id 0 and no password
var (
    redis *godis.Client
    config *Config
)



// function to resolve a shorturl and redirect
func resolve(ctx *web.Context, short string) {
    redirect, err := redis.Get(short)
    if err == nil {
        ctx.Redirect(301, redirect.String())
    } else {
        ctx.Redirect(301, "https://www.youtube.com/watch?v=jRHmvy5eaG4")
    }
}



func isValidUrl(rawurl string) (u *url.URL, err os.Error){
    if len(rawurl) == 0{
        return nil, os.NewError("empty url")
    }
    // XXX this needs some love...
    if !strings.HasPrefix(rawurl, HTTP){
        rawurl = fmt.Sprintf("%s://%s", HTTP, rawurl)
    }
    return url.Parse(rawurl)
}


// function to shorten and store a url
func shorten(ctx *web.Context, data string){
    const(
        jsntmpl = "{\"url\" : \"%s\", \"longurl\" : \"%s\"}\n"
    )
    host := config.GetStringDefault("hostname", "localhost")
    r, _ := ctx.Request.Params["url"]
    theUrl, err := isValidUrl(string(r))
    if err == nil{
        ctr, _ := redis.Incr(COUNTER)
        encoded := Encode(ctr)
        // fire and forget
        go redis.Set(encoded, theUrl.Raw)

        ctx.SetHeader("Content-Type", "application/json", true)
        location := fmt.Sprintf("%s://%s/%s", HTTP, host, encoded)
        ctx.SetHeader("Location", location, true)
        ctx.StartResponse(201)
        ctx.WriteString(fmt.Sprintf(jsntmpl, location, theUrl.Raw))
    }else{
       ctx.Redirect(404, "/")
    }
}

func bootstrap(path string) os.Error {
    config = NewConfig(path)
    config.Parse()
    host := config.GetStringDefault("redis.address", "tcp:localhost:6379")
    db := config.GetIntDefault("redis.database", 0)
    passwd := config.GetStringDefault("redis.password", "")

    redis = godis.New(host, db, passwd)
    return nil
}




// main function that inits the routes in web.go
func main() {
    flag.Parse()
    cfgFile := flag.Arg(0)
    err := bootstrap(cfgFile)
    if err == nil {
        // this could go to bootstrap as well
        web.Post("/shorten/(.*)", shorten)
        web.Get("/(.*)", resolve)
        listen := config.GetStringDefault("listen", "0.0.0.0")
        port := config.GetStringDefault("port", "9999")
        web.Run(fmt.Sprintf("%s:%s", listen, port))
    }
}

