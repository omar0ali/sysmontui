package screentui

// each element added on screen should be a point
// a point is a place where we can start to render something
// i.e we can draw box, add a text, center it
// each entity can have multiple points

type Point struct {
	X, Y float64
}

func P(X, Y float64) Point {
	return Point{X, Y}
}
