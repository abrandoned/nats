package main

/*
#include <stdlib.h>
*/
import "C"
import "fmt"
import "unsafe"
import "math/rand"
import "strconv"
import "time"
import "github.com/nats-io/nats"
import "github.com/streamrail/concurrent-map"

var connectionMap = cmap.New()

//export Close
func Close(connection_id int) {
  connection_id_string := strconv.Itoa(connection_id)
	nc, _ := connectionMap.Get(connection_id_string)
  connection := nc.(*nats.Conn)
	connection.Close()
	connectionMap.Remove(connection_id_string)
}

//export CloseAll
func CloseAll() {
	for item := range connectionMap.Iter() {
		item.Val.(*nats.Conn).Close()
	}

	connectionMap = cmap.New()
}

//export Connect
func Connect(c_url *C.char) int32 {
	url := C.GoString(c_url)
	rando := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand_int := rando.Int31()

	nc, _ := nats.Connect(url)
  for {
    connection_id_string := strconv.FormatInt(int64(rand_int), 10)
    if connectionMap.Has(connection_id_string) == true {
      rand_int = rando.Int31()
    } else {
	    connectionMap.Set(connection_id_string, nc)
      break
    }
  }

	return rand_int
}

//export Flush
func Flush(connection_id int) {
  connection_id_string := strconv.Itoa(connection_id)
	nc, _ := connectionMap.Get(connection_id_string)
	nc.(*nats.Conn).Flush()
}

//export FlushAll
func FlushAll() {
	for item := range connectionMap.Iter() {
		item.Val.(*nats.Conn).Flush()
	}
}

//export Publish
func Publish(connection_id int, c_subject *C.char, c_message *C.char, message_length C.int) {
  connection_id_string := strconv.Itoa(connection_id)
	nc, _ := connectionMap.Get(connection_id_string)
	nc.(*nats.Conn).Publish(C.GoString(c_subject), C.GoBytes(unsafe.Pointer(c_message), message_length))
}

//export Request
func Request(connection_id int, c_subject *C.char, c_message *C.char) *C.char {
  connection_id_string := strconv.Itoa(connection_id)
	subject := C.GoString(c_subject)
	message := C.GoString(c_message)

	nc, _ := connectionMap.Get(connection_id_string)
	msg, err := nc.(*nats.Conn).Request(subject, []byte(message), 10*time.Millisecond)

	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}
	// TODO handle error responses (Msg struct is empty after it)

	response_string := C.CString(string(msg.Subject))
	defer C.free(unsafe.Pointer(response_string))

	return response_string
}

func main() {}
