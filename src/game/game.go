package game

type Board struct {
	Height, Width int
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Food struct {
	Loc Point `json:"location"`
}
