package main

import (
	"encoding/json"
	"errors"
    "io"
    "os"
    "path"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
    "code.google.com/p/gorilla/mux"
	"github.com/simonz05/godis"
)

const (
	// special key in redis, that is our global counter
	COUNTER = "__counter__"
	HTTP    = "http"
	ROLL    = "https://www.youtube.com/watch?v=jRHmvy5eaG4"
)

var (
	redis  *godis.Client
	config *Config
)

type KurzUrl struct {
	Key          string
	ShortUrl     string
	LongUrl      string
	CreationDate int64
	Clicks       int64
}

// Converts the KurzUrl to JSON.
func (k KurzUrl) Json() []byte {
	b, _ := json.Marshal(k)
	return b
}

// Creates a new KurzUrl instance. The Given key, shorturl and longurl will
// be used. Clicks will be set to 0 and CreationDate to time.Nanoseconds()
func NewKurzUrl(key, shorturl, longurl string) *KurzUrl {
	kurl := new(KurzUrl)
	kurl.CreationDate = time.Now().UnixNano()
	kurl.Key = key
	kurl.LongUrl = longurl
	kurl.ShortUrl = shorturl
	kurl.Clicks = 0
	return kurl
}

// stores a new KurzUrl for the given key, shorturl and longurl. Existing
// ones with the same url will be overwritten
func store(key, shorturl, longurl string) *KurzUrl {
	kurl := NewKurzUrl(key, shorturl, longurl)
	go redis.Hset(kurl.Key, "LongUrl", kurl.LongUrl)
	go redis.Hset(kurl.Key, "ShortUrl", kurl.ShortUrl)
	go redis.Hset(kurl.Key, "CreationDate", kurl.CreationDate)
	go redis.Hset(kurl.Key, "Clicks", kurl.Clicks)
	return kurl
}

// loads a KurzUrl instance for the given key. If the key is
// not found, os.Error is returned.
func load(key string) (*KurzUrl, error) {
	if ok, _ := redis.Hexists(key, "ShortUrl"); ok {
		kurl := new(KurzUrl)
		kurl.Key = key
		reply, _ := redis.Hmget(key, "LongUrl", "ShortUrl", "CreationDate", "Clicks")
		kurl.LongUrl, kurl.ShortUrl, kurl.CreationDate, kurl.Clicks =
			reply.Elems[0].Elem.String(), reply.Elems[1].Elem.String(),
			reply.Elems[2].Elem.Int64(), reply.Elems[3].Elem.Int64()
		return kurl, nil
	}
	return nil, errors.New("unknown key: " + key)
}

func fileExists(dir string) bool {
    info, err := os.Stat(dir)
    if err != nil {
        return false
    }

    return !info.IsDir()
}


// function to display the info about a KurzUrl given by it's Key
func info(w http.ResponseWriter, r *http.Request){
    short := mux.Vars(r)["short"]
    if strings.HasSuffix(short, "+"){
        short = strings.Replace(short, "+", "", 1)
    }

	kurl, err := load(short)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(kurl.Json())
		io.WriteString(w, "\n")
	} else {
		http.Redirect(w, r, ROLL, http.StatusNotFound)
	}
}

// function to resolve a shorturl and redirect
func resolve(w http.ResponseWriter, r *http.Request){

    short := mux.Vars(r)["short"]
	kurl, err := load(short)
	if err == nil {
		go redis.Hincrby(kurl.Key, "Clicks", 1)
		http.Redirect(w, r, kurl.LongUrl, http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, ROLL, http.StatusMovedPermanently)
	}
}

// Determines if the string rawurl is a valid URL to be stored.
func isValidUrl(rawurl string) (u *url.URL, err error) {
	if len(rawurl) == 0 {
		return nil, errors.New("empty url")
	}
	// XXX this needs some love...
	if !strings.HasPrefix(rawurl, HTTP) {
		rawurl = fmt.Sprintf("%s://%s", HTTP, rawurl)
	}
	return url.Parse(rawurl)
}

// function to shorten and store a url
func shorten(w http.ResponseWriter, r *http.Request){
	host := config.GetStringDefault("hostname", "localhost")
    var leUrl string
    if (r.Method == "GET"){
        leUrl = mux.Vars(r)["url"]
    }else{
        leUrl = r.FormValue("url")
    }
	theUrl, err := isValidUrl(string(leUrl))
	if err == nil {
		ctr, _ := redis.Incr(COUNTER)
		encoded := Encode(ctr)
		location := fmt.Sprintf("%s://%s/%s", HTTP, host, encoded)
		store(encoded, location, theUrl.String())
		// redirect to the info page
		http.Redirect(w, r,  location + "+", http.StatusMovedPermanently)
	} else {
		http.Redirect(w, r, ROLL, http.StatusNotFound)
	}
}

//Returns a json array with information about the last shortened urls. If data 
// is a valid integer, that's the amount of data it will return, otherwise
// a maximum of 10 entries will be returned.
func latest(w http.ResponseWriter, r *http.Request){
    data := mux.Vars(r)["data"]
	howmany, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		howmany = 10
	}
	c, _ := redis.Get(COUNTER)

	last := c.Int64()
	upTo := (last - howmany)

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, "{ \"urls\" : [")
	for i := last; i > upTo && i > 0; i -= 1 {
		kurl, err := load(Encode(i))
		if err == nil {
			w.Write(kurl.Json())
			if i != upTo+1 {
				io.WriteString(w, ",")
			}
		}
	}
	io.WriteString(w, "] }")
	io.WriteString(w, "\n")
}


func static(w http.ResponseWriter, r *http.Request){
    fname := mux.Vars(r)["fileName"]
    // empty means, we want ot serve the index file. Due to a bug in http.serveFile
    // the file cannot be called index.html, anything else is fine.
    if fname == ""{
        fname = "index.htm"
    }
    staticDir := config.GetStringDefault("static-directory", "")
    staticFile := path.Join(staticDir, fname)
    if fileExists(staticFile){
        http.ServeFile(w, r, staticFile)
    }
}


func main() {
    flag.Parse()
	path := flag.Arg(0)

	config = NewConfig(path)
	config.Parse()

	host := config.GetStringDefault("redis.address", "tcp:localhost:6379")
	db := config.GetIntDefault("redis.database", 0)
	passwd := config.GetStringDefault("redis.password", "")

	redis = godis.New(host, db, passwd)

    router := mux.NewRouter()
	router.HandleFunc("/shorten/{url:(.*$)}", shorten)

    router.HandleFunc("/{short:([a-zA-Z0-9]+$)}", resolve)
	router.HandleFunc("/{short:([a-zA-Z0-9]+)\\+$}", info)
	router.HandleFunc("/info/{short:[a-zA-Z0-9]+}", info)
	router.HandleFunc("/latest/{data:[0-9]+}", latest)

	router.HandleFunc("/{fileName:(.*$)}", static)



    listen := config.GetStringDefault("listen", "0.0.0.0")
    port := config.GetStringDefault("port", "9999")
    s := &http.Server{
        Addr:           listen + ":" + port,
        Handler:        router,
    }
    s.ListenAndServe()
}
