package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "sync"
  //"encoding/json"

  "github.com/julienschmidt/httprouter"
  "github.com/rs/cors"
)

// create a todo struct
type todo struct {
  id string
  title string
  completed string
}

// create a store struct
type store struct {
  data map[string]todo
  m sync.RWMutex
}

// set local vars
var (
  addr = flag.String("addr", ":8080", "http service address")
  s = store{
    data: map[string]todo{},
    m: sync.RWMutex{},
  }
)

func main(){
  // begin by parsing the command line
  flag.Parse()

  // get and set router
  // then set the routes and actions
  router := httprouter.New()
  router.GET("/todo/:key", show)
  router.GET("/todo", show)
  router.PUT("/todo/:key/:title/:completed", add)
  router.PUT("/update/todo/:key/:title/:completed", update)
  router.DELETE("/todo/:key", remove)

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

    // begin: convert data to json string and set to result
    s.m.RLock()
    // get data
    m := s.data

    // set index and sizeOfStore
    i := 0
    sizeOfStore := len(m)

    // iterate through data map, make values into json string, save in result
    result := "{\"todos\": ["
    for key := range m {
      sp := s.data[key]
      if i < sizeOfStore-1 {
        result = result + "{\"id\":\""+sp.id+"\", \"title\":\""+sp.title+"\", \"completed\":\""+sp.completed+"\"},"
      } else {
        result = result + "{\"id\":\""+sp.id+"\", \"title\":\""+sp.title+"\", \"completed\":\""+sp.completed+"\"}"
      }

      i++
    }
    result = result + "]}"
    s.m.RUnlock()
    // end: convert data to json string and set to result

    // return result to caller
    fmt.Fprintf(w, "%s", result)

    return;
  }

  // get key=data pair, make json string, and return result to caller
  s.m.RLock()
  sp := s.data[k]
  fmt.Fprintf(w, "{\"id\":\"%s\", \"title\":\"%s\", \"completed\":\"%s\"}", sp.id, sp.title, sp.completed)
  s.m.RUnlock()

  return;
}

func add(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  // get and set key, title, completed, and testTitle
  k := p.ByName("key")
  title := p.ByName("title")
  completed := p.ByName("completed")

  myTodo := todo{
    id: p.ByName("key"),
    title: p.ByName("title"),
    completed: p.ByName("completed"),
  }

  // add data in the store with key value
  s.m.RLock()
  s.data[k] = myTodo
  s.m.RUnlock()

  // return result
  fmt.Fprintf(w, "Added: data[%s] = %s, completed:%s", k, title, completed)

  return;
}

func update(w http.ResponseWriter, r *http.Request, p httprouter.Params){
  // get and set key, title, completed, and testTitle
  k := p.ByName("key")
  title := p.ByName("title")
  completed := p.ByName("completed")

  myTodo := s.data[k]

  // update data in the store with key value
  s.m.RLock()
  if title != "" {
    myTodo.title = title;
  }

  myTodo.completed = completed;

  s.data[k] = myTodo
  s.m.RUnlock()

  return;
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

  return;
}

// this func used for developer reference only
func array_map_iterators(){
  arr := []string{"a", "b", "c"}

  for index, value := range arr {
    fmt.Println("index:", index, "value:", value)
  }

  m := make(map[string]string)
  m["a"] = "alpha"
  m["b"] = "beta"

  for key, value := range m {
    fmt.Println("key:", key, "value:", value)
  }

  return;
}
