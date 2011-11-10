package main

import (
    "github.com/hoisie/web.go"
    "encoding/base64"
    "encoding/binary"
    "strings"
    /*"fmt"*/
)

var counter = 0
var allurls = map[string]string{
}

func resolve(ctx *web.Context, short string) {
    redirect, found := allurls[short]
    //fmt.Printf(encoded)
    if found == true {
        ctx.Redirect(302, redirect)
    } else{
        ctx.Redirect(404, "")
    }
}

func store(ctx *web.Context, url string) string {
    if ! strings.HasPrefix(url, "http"){
        url = "http://" + url
    }
    counter += 1
    request := ctx.Request
    b := make([]byte, 4)
    binary.BigEndian.PutUint32(b[0:], uint32(counter))
    encoded := base64.StdEncoding.EncodeToString(b)
    allurls[encoded] = url
    return request.Proto + request.Host + "/" + encoded

}


func main() {
    web.Get("/store/(.*)", store)
    web.Get("/(.*)", resolve)
    web.Run("0.0.0.0:9999")
}

