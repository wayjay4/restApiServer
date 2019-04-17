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

// create a store struct
type store struct {
  data map[string]string
  m sync.RWMutex
}

// set local vars
var (
  addr = flag.String("addr", ":8080", "http service address")
  s = store{
    data: map[string]string{},
    m: sync.RWMutex{},
  }
)

func main(){
  // begin by parsing the command line
  flag.Parse()

  // get and set router
  // then set the routes and actions
  router := httprouter.New()
  router.GET("/entry/:key", show)
  router.GET("/list", show)
  router.PUT("/entry/:key/:title/:completed", update)
  router.DELETE("/entry/:key", remove)

  // get and set CORS
  // then set options to allow cross origin and allowed methods for server
  c := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
    AllowCredentials: true,
    // Enable Debugging for testing, consider disabling in production
    Debug: true,
  })

  // set CORS as handler for router
  handler := c.Handler(router)

  // start server and capture any errors
  err := http.ListenAndServe(*addr, handler)

  // if there were any errors, then log the error
  if err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}

func show(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  // get and set key
  k := p.ByName("key")

  // if there is no key, then return 'list' of contents of store
  if k == "" {

    // convert data to json string and set to result
    s.m.RLock()
    result, _ := json.Marshal(s.data)
    s.m.RUnlock()

    // return result to caller
    fmt.Fprintf(w, "%s", result)

    return
  }

  // get key=data pair and return result to caller
  s.m.RLock()
  fmt.Fprintf(w, "Read entry: data[%s] = %s", k, s.data[k])
  s.m.RUnlock()

  return
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  // get and set key, title, completed, and testTitle
  k := p.ByName("key")
  title := p.ByName("title")
  completed := p.ByName("completed")
  testTitle := p.ByName("testTitle")

  // update data in the store with key value
  s.m.RLock()
  s.data[k] = title
  s.m.RUnlock()

  // return result
  fmt.Fprintf(w, "Updated: data[%s] = %s, completed:%s [testTitle: %s]", k, title, completed, testTitle)
}

func remove(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  // get and set key
  k := p.ByName("key")

  // delete data in the store with the key value
  s.m.RLock()
  v := s.data[k]
  delete(s.data, k);
  s.m.RUnlock()

  // return result
  fmt.Fprintf(w, "Deleted: data[%s] = %s", k, v)
}
