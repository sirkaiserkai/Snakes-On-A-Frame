package world

import (
	"encoding/json"
	"events"
	"game"
	"log"
	"math/rand"
	"snake"
	"sync"
)

type World struct {
	Mutex  sync.Mutex
	board  game.Board
	food   game.Food
	snakes []*snake.Snake
	events []string
	Ed     *events.EventDispatcher
}

func randomPoint(min, max int) game.Point {
	x := rand.Intn(max+min) + min
	y := rand.Intn(max+min) + min

	return game.Point{X: x, Y: y}
}

func NewWorld(ed *events.EventDispatcher, width, height int) World {
	b := game.Board{Width: width, Height: height}
	w := World{
		Mutex:  sync.Mutex{},
		board:  b,
		snakes: []*snake.Snake{},
		food:   game.Food{Loc: randomPoint(0, width)},
		events: []string{},
		Ed:     ed,
	}

	return w
}

func (w *World) NumberOfSnakes() int {
	w.Mutex.Lock()
	numberOfSnakes := len(w.snakes)
	w.Mutex.Unlock()

	return numberOfSnakes
}

func (w *World) AddSnake(s *snake.Snake) {
	// TODO: Make thread safe
	w.Mutex.Lock()
	w.snakes = append(w.snakes, s)
	w.Ed.Listeners = append(w.Ed.Listeners, s)
	w.Mutex.Unlock()
}

func (w *World) RemoveSnake(s0 *snake.Snake) {
	w.Mutex.Lock()
	for i, s1 := range w.snakes {
		if s0.ID == s1.ID {
			w.snakes = append(w.snakes[:i], w.snakes[i+1:]...)

			w.Ed.Listeners = []events.Listener{}
			for _, s := range w.snakes {
				w.Ed.Listeners = append(w.Ed.Listeners, s)
			}
		}
	}
	w.Mutex.Unlock()
}

func (w *World) Step() {
	w.Mutex.Lock()
	event := events.Event{EventType: "step"}
	w.Ed.AddEvent(event)
	w.Mutex.Unlock()
}

func (w *World) CalculateWorldState() {
	w.Mutex.Lock()

	// Mutex to prevent snakesToRemove dict being modified simultaneously
	dicMutex := sync.Mutex{}
	snakesToRemove := make(map[string]snake.Snake)

	// group which waits until all calculations are completed.
	wg := sync.WaitGroup{}

	for _, s0 := range w.snakes {

		wg.Add(1)
		go func(s snake.Snake, b game.Board) {
			defer wg.Done()
			if s.WallCollision(b) {
				dicMutex.Lock()
				snakesToRemove[s.ID] = s
				dicMutex.Unlock()
			}
		}(*s0, w.board)

		wg.Add(1)
		go func(s snake.Snake) {
			defer wg.Done()
			for i, bodyPart := range s.Body {
				if i != 0 {
					if s.Head() == bodyPart {
						dicMutex.Lock()
						snakesToRemove[s.ID] = s
						dicMutex.Unlock()
					}
				}
			}
		}(*s0)

		for _, s1 := range w.snakes {
			wg.Add(1)
			go func(s0, s1 snake.Snake) {
				defer wg.Done()
				if s0.ID != s1.ID {
					if s0.SnakeCollision(s1) {
						dicMutex.Lock()
						snakesToRemove[s0.ID] = s0
						snakesToRemove[s1.ID] = s1
						dicMutex.Unlock()
					}
				}
			}(*s0, *s1)
		}

		if s0.FoodCollision(w.food) {
			s0.Score += 10
			s0.AddBodyPart()
			w.food.Loc = randomPoint(0, w.board.Width)

		}
	}

	wg.Wait()
	snakes := w.snakes[:0]
	w.Ed.Listeners = []events.Listener{}
	for _, s := range w.snakes {
		if _, ok := snakesToRemove[s.ID]; !ok {
			snakes = append(snakes, s)
			w.Ed.Listeners = append(w.Ed.Listeners, s)
		}
	}

	w.snakes = snakes
	w.Mutex.Unlock()
}

func (w *World) GetState() string {

	ss, err := json.Marshal(w.snakes)
	if err != nil {
		log.Println("Error in GetState: %s", err)
		return "N/A"
	}

	m := make(map[string]interface{})
	m["snakes"] = string(ss)

	food, err := json.Marshal(w.food)
	if err != nil {
		log.Println("Error in GetState: %s", err)
		return "N/A"
	}

	m["food"] = string(food)

	msg, err := json.Marshal(m)
	if err != nil {
		log.Println("Error in GetState: %s", err)
		return "N/A"
	}

	return string(msg)
}
