package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "sync"
  "encoding/json"

  "github.com/julienschmidt/httprouter"
  "github.com/rs/cors"
)

type store struct {
  data map[string]string
  m sync.RWMutex
}

var (
  addr = flag.String("addr", ":8080", "http service address")
  s = store{
    data: map[string]string{},
    m: sync.RWMutex{},
  }
)

func main(){
  flag.Parse()
  router := httprouter.New()
  router.GET("/entry/:key", show)
  router.GET("/list", show)
  router.PUT("/entry/:key/:title/:completed", update)
  router.DELETE("/entry/:key", remove)

  c := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
    AllowCredentials: true,
    // Enable Debugging for testing, consider disabling in production
    Debug: true,
})

  handler := c.Handler(router)
  err := http.ListenAndServe(*addr, handler)
  //err := http.ListenAndServe(*addr, router)

  if err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  k := p.ByName("key")
  if k == "" {
    s.m.RLock()
    //fmt.Fprintf(w, "Read list: %v", s.data)

    res1B, _ := json.Marshal(s.data)
    fmt.Fprintf(w, "%s", res1B)

    s.m.RUnlock()
    return
  }

  s.m.RLock()
  fmt.Fprintf(w, "Read entry: data[%s] = %s", k, s.data[k])
  s.m.RUnlock()
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  k := p.ByName("key")
  title := p.ByName("title")
  completed := p.ByName("completed")
  testTitle := p.ByName("testTitle")

  s.m.RLock()
  s.data[k] = title
  s.m.RUnlock()

  fmt.Fprintf(w, "Updated: data[%s] = %s, completed:%s [testTitle: %s]", k, title, completed, testTitle)
}

func remove(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  k := p.ByName("key")

  s.m.RLock()
  v := s.data[k]
  delete(s.data, k);
  s.m.RUnlock()

  fmt.Fprintf(w, "Deleted: data[%s] = %s", k, v)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
  (*w).Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token")
}
