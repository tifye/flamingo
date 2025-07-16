package render

type Component interface {
	Name() string
	Attributes() map[string]any
	SetAttribute(key string, val any)
	Children() []Component
}

type Renderer interface {
	Render(comps ...Component)
	Append(parent Component, child Component)
	NewComponent(name string) Component
}
