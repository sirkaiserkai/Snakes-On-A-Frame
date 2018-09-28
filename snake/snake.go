package snake

import (
	"crypto/rand"
	"github.com/sirkaiserkai/Snakes-On-A-Frame/events"
	"fmt"
	"github.com/sirkaiserkai/Snakes-On-A-Frame/game"
	"log"
	mRand "math/rand"
	"sync"

	"github.com/googollee/go-socket.io"
)

var colors = []string{
	"LightGreen",
	"LightSkyBlue",
	"Pink",
	"Navy",
	"Teal",
	"SteelBlue",
	"RebeccaPurple",
	"Coral",
	"DarkRed",
	"MediumPurple",
}

func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

type Direction int

const (
	Up    Direction = iota
	Down  Direction = iota
	Right Direction = iota
	Left  Direction = iota
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Right:
		return "Right"
	case Left:
		return "Left"
	default:
		return "Unknown"
	}
}

func GetDirection(direction string) Direction {
	switch direction {
	case "up":
		return Up
	case "down":
		return Down
	case "right":
		return Right
	case "left":
		return Left
	}

	return Right
}

type Snake struct {
	m         sync.Mutex
	ID        string       `json:"id"`
	Username  string       `json:"username"`
	Body      []game.Point `json:"body"`
	direction Direction    `json:"direction"`
	Color     string       `json:"color"`
	Score     int          `json:"score"`

	addBodyPart bool
	so          *socketio.Socket
}

func NewSnake(so *socketio.Socket) Snake {
	snake := Snake{}

	snake.m = sync.Mutex{}
	snake.direction = Right
	snake.Color = colors[mRand.Intn(len(colors))]
	snake.ID = pseudo_uuid()
	snake.so = so
	snake.addBodyPart = false
	snake.Body = []game.Point{
		game.Point{X: 9, Y: 0},
		game.Point{X: 8, Y: 0},
		game.Point{X: 7, Y: 0},
		game.Point{X: 6, Y: 0},
		game.Point{X: 4, Y: 0},
		game.Point{X: 3, Y: 0},
		game.Point{X: 2, Y: 0},
		game.Point{X: 1, Y: 0}}

	return snake
}

func (s *Snake) Length() int {
	return len(s.Body)
}

func (s *Snake) Head() game.Point {
	return s.Body[0]
}

func (s *Snake) NextStep() {
	nextPoint := s.Body[0]
	switch s.direction {
	case Up:
		nextPoint.Y--
	case Down:
		nextPoint.Y++
	case Right:
		nextPoint.X++
	case Left:
		nextPoint.X--
	}

	modifier := 1
	if s.addBodyPart {
		modifier = 0
		s.addBodyPart = false
	}
	s.Body = append([]game.Point{nextPoint}, s.Body[:len(s.Body)-modifier]...)
}

func (s *Snake) AddBodyPart() {
	s.addBodyPart = true
}

func (s *Snake) SetDirection(d Direction) {
	switch d {
	case Up:
		if s.direction != Down {
			s.direction = Up
		}
	case Down:
		if s.direction != Up {
			s.direction = Down
		}
	case Right:
		if s.direction != Left {
			s.direction = Right
		}
	case Left:
		if s.direction != Right {
			s.direction = Left
		}
	default:
		return
	}
}

func (s *Snake) FoodCollision(f game.Food) bool {
	return s.Head().X == f.Loc.X && s.Head().Y == f.Loc.Y
}

func (s *Snake) WallCollision(b game.Board) bool {
	x := s.Head().X
	y := s.Head().Y

	return (x > b.Width ||
		x < 0 ||
		y > b.Height ||
		y < 0)
}

func (s0 *Snake) SnakeCollision(s1 Snake) bool {
	/*for _, i := range s0.Body {
		for _, j := range s1.Body {
			if i == j {
				return true // collis
			}
		}
	}*/

	for _, b := range s1.Body {
		if s0.Head() == b {
			return true // collision
		}
	}

	return false // No collisions
}

func (s *Snake) HandleEvent(event events.Event, wg *sync.WaitGroup) {
	log.Printf("snake: %s handling event: %s", s.ID, event.EventType)
	switch event.EventType {
	case "move":
		s.HandleMoveEvent(event)
	case "step":
		s.HandleStepEvent(event)
	}
	wg.Done()
}

func (s *Snake) HandleMoveEvent(event events.Event) {
	s.HandleMoveDict("")
}

/*func (s *Snake) HandleMoveDict(moveDict map[string]interface{}) {
	id, ok := moveDict["id"].(string)
	if !ok {
		return
	}

	if s.ID != id {
		return
	}

	direction, ok := moveDict["direction"].(string)
	if !ok {
		return
	}

	d := GetDirection(direction)
	s.SetDirection(d)
}*/

func (s *Snake) HandleMoveDict(direction string) {
	/*id, ok := moveDict["id"].(string)
	if !ok {
		return
	}

	if s.ID != id {
		return
	}

	direction, ok := moveDict["direction"].(string)
	if !ok {
		return
	}*/

	d := GetDirection(direction)
	s.SetDirection(d)
}

func (s *Snake) HandleStepEvent(event events.Event) {
	s.NextStep()
}
