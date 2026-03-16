package screentui

type Entity interface {
	Update(d float64)
	Render()
	Events()
}
