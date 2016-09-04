package main

// #include <stdio.h>
// #include <stdlib.h>

import "C"
import "unsafe"
import "math/rand"
import "time"
import "github.com/nats-io/nats"

var connectionMap = make(map[int64]*nats.Conn)

//export Close
func Close(connection int64) {
  nc := connectionMap[connection]
  nc.Close()
}

//export CloseAll
func CloseAll() {
  for _, value := range connectionMap {
    value.Close()
  }

  connectionMap = make(map[int64]*nats.Conn)
}

//export Connect
func Connect(c_url *C.char) (int64)  {
  url := C.GoString(c_url)
  rando := rand.NewSource(time.Now().UnixNano())
  rand64 := rando.Int63()

  nc, _ := nats.Connect(url)
  connectionMap[rand64] = nc
  return rand64
}

//export Flush
func Flush(connection int64) {
  nc := connectionMap[connection]
  nc.Flush()
}

//export FlushAll
func FlushAll() {
  for _, value := range connectionMap {
    value.Flush()
  }
}

//export Publish
func Publish (connection int64, c_subject *C.char, c_message *C.char) {
  subject := C.GoString(c_subject)
  message := C.GoString(c_message)
  nc := connectionMap[connection]
  nc.Publish(subject, []byte(message))
}

//export Request
func Request (connection int64, c_subject *C.char, c_message *C.char) *C.char {
  subject := C.GoString(c_subject)
  message := C.GoString(c_message)
  nc := connectionMap[connection]
  msg, _ := nc.Request(subject, []byte(message), 10 * time.Millisecond)

  response_string := C.CString(msg.Data, len(msg.Data))
  defer C.free(unsafe.Pointer(response_string))

  return response_string
}

func main() {}
