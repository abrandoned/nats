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
func Close(connection int64) {
	nc, _ := connectionMap.Get(strconv.FormatInt(connection, 10))
	nc.Close()
	connectionMap.Remove(strconv.FormatInt(connection, 10))
}

//export CloseAll
func CloseAll() {
	for item := range connectionMap.Iter() {
		item.Val.Close()
	}

	connectionMap = cmap.New()
}

//export Connect
func Connect(c_url *C.char) int64 {
	url := C.GoString(c_url)
	rando := rand.NewSource(time.Now().UnixNano())
	rand64 := rando.Int63()

	nc, _ := nats.Connect(url)
	connectionMap.Set(strconv.FormatInt(rand64, 10), nc)
	return rand64
}

//export Flush
func Flush(connection int64) {
	nc, _ := connectionMap.Get(strconv.FormatInt(connection, 10))
	nc.Flush()
}

//export FlushAll
func FlushAll() {
	for item := range connectionMap.Iter() {
		item.Val.Flush()
	}
}

//export Publish
func Publish(connection int64, c_subject *C.char, c_message *C.char) {
	subject := C.GoString(c_subject)
	message := C.GoString(c_message)
	nc, _ := connectionMap.Get(strconv.FormatInt(connection, 10))
	nc.Publish(subject, []byte(message))
}

//export Request
func Request(connection int64, c_subject *C.char, c_message *C.char) *C.char {
	subject := C.GoString(c_subject)
	message := C.GoString(c_message)

	nc, _ := connectionMap.Get(strconv.FormatInt(connection, 10))
	msg, err := nc.Request(subject, []byte(message), 10*time.Millisecond)

	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
	}
	// TODO handle error responses (Msg struct is empty after it)

	response_string := C.CString(string(msg.Subject))
	defer C.free(unsafe.Pointer(response_string))

	return response_string
}

func main() {}
