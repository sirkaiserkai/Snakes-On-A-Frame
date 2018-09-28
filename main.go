package main

import (
	"github.com/sirkaiserkai/Snakes-On-A-Frame/events"
	"log"
	"net/http"
	"github.com/sirkaiserkai/Snakes-On-A-Frame/snake"
	"time"
	"github.com/sirkaiserkai/Snakes-On-A-Frame/world"

	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
)

const port = ":8000"
const stepTimeInterval = 60 * time.Millisecond

// TODO: Set up snake game
// - Create snake struct
// - set up game board struct
// - set up handler methods for socket events
// - set up player type

var ed events.EventDispatcher = events.EventDispatcher{}
var w world.World = world.NewWorld(&ed, 45, 45)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join("game")

		snake := snake.NewSnake(&so)
		log.Println("Adding snake")
		w.AddSnake(&snake)

		// Have snake subscribe to event dispatcher

		so.On("move", snake.HandleMoveDict)

		so.On("disconnection", func() {
			log.Println("snake: ", snake.ID, "disconnected")
			w.RemoveSnake(&snake)
		})
	})

	// Main game loop
	go func(server *socketio.Server) {
		stateTimer := time.NewTimer(stepTimeInterval)
	GameLoop:
		for {

			if w.NumberOfSnakes() == 0 {
				continue GameLoop
			}

			select {
			case <-stateTimer.C: // Send new state to clients

				//log.Println("Sending state")
				w.Step()
				w.CalculateWorldState()
				log.Println("Snakes: ", w.GetState())
				server.BroadcastTo("game", "new_state", w.GetState())
				stateTimer.Reset(stepTimeInterval)
				//log.Println("Timer reset.")
				// Process next step.
			default:
				// log.Println("Nothing")
			}
		}
	}(server)

	r := mux.NewRouter()

	http.Handle("/socket.io/", server)
	http.Handle("/", r)
	// Serves HTML documents
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./public"))))

	log.Printf("Serving at localhost:%s...\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
