package main

import "log"
import "net/http"

import "github.com/googollee/go-socket.io"

// Message json object
type Message struct {
  Username string `json:"username"`
  Message  string `json:"message"`
}

func main() {
  server, err := socketio.NewServer(nil)

  if err != nil {
    log.Fatal(err)
  }

  server.On("connection", func(socket socketio.Socket) {
    log.Println("socket " + socket.Id() + " connected")

    socket.Join("chat")

    socket.On("chatMessage", func(msg Message) {
      log.Println("emit: " + msg.Message)
      // Send to self
      socket.Emit("chatMessage", msg)
      // Send to everyone else
      socket.BroadcastTo("chat", "chatMessage", msg)
    })

    socket.On("disconnect", func() {
      log.Println("on disconnect")
    })
  })

  server.On("error", func(socket socketio.Socket, err error) {
    log.Println("error:", err)
  })

  http.Handle("/socket.io/", server)
  http.Handle("/", http.FileServer(http.Dir("./public")))
  log.Println("server running on localhost:8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
