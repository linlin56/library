// Package uml — UseCaseDiagram builds a PlantUML use case diagram
// fluently, in the same way as SequenceDiagram.
package uml

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Diagram is an interface shared by every PlantUML diagram type.
// It allows different types (sequence, use case…) to be stored
// in the same map and written with a single loop.
type Diagram interface {
	WriteFile(path string) error
}

// UseCaseDiagram builds a PlantUML use case diagram.
type UseCaseDiagram struct {
	title string
	lines []string
}

// NewUseCaseDiagram creates a new use case diagram with a title.
func NewUseCaseDiagram(title string) *UseCaseDiagram {
	return &UseCaseDiagram{title: title}
}

func (u *UseCaseDiagram) add(l string) *UseCaseDiagram {
	u.lines = append(u.lines, l)
	return u
}

// Direction sets the diagram orientation (e.g. "left to right direction").
func (u *UseCaseDiagram) Direction(dir string) *UseCaseDiagram {
	return u.add(dir)
}

// Actor declares an actor with an alias and a displayed label.
func (u *UseCaseDiagram) Actor(alias, label string) *UseCaseDiagram {
	return u.add(fmt.Sprintf(`actor "%s" as %s`, label, alias))
}

// ActorInherits declares an inheritance relation between actors (child inherits from parent).
// Produces: parent <|-- child
func (u *UseCaseDiagram) ActorInherits(child, parent string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s <|-- %s", parent, child))
}

// UseCase declares a use case with an alias and a label.
func (u *UseCaseDiagram) UseCase(alias, label string) *UseCaseDiagram {
	return u.add(fmt.Sprintf(`usecase "%s" as %s`, label, alias))
}

// RectangleStart opens a rectangle block (system boundary).
func (u *UseCaseDiagram) RectangleStart(label string) *UseCaseDiagram {
	return u.add(fmt.Sprintf(`rectangle "%s" {`, label))
}

// RectangleEnd closes a rectangle block.
func (u *UseCaseDiagram) RectangleEnd() *UseCaseDiagram {
	return u.add("}")
}

// Association links an actor to a use case.
func (u *UseCaseDiagram) Association(from, to string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s --> %s", from, to))
}

// Include links two use cases with the <<include>> relation.
func (u *UseCaseDiagram) Include(from, to string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s ..> %s : <<include>>", from, to))
}

// Extend links two use cases with the <<extend>> relation.
func (u *UseCaseDiagram) Extend(from, to string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s ..> %s : <<extend>>", from, to))
}

// Blank inserts an empty line for readability.
func (u *UseCaseDiagram) Blank() *UseCaseDiagram {
	return u.add("")
}

// String returns the full diagram in PlantUML syntax.
func (u *UseCaseDiagram) String() string {
	var b strings.Builder
	b.WriteString("@startuml\n")
	if u.title != "" {
		fmt.Fprintf(&b, "title %s\n\n", u.title)
	}
	for _, l := range u.lines {
		b.WriteString(l + "\n")
	}
	b.WriteString("@enduml\n")
	return b.String()
}

// WriteTo writes the diagram to an io.Writer.
func (u *UseCaseDiagram) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, u.String())
	return int64(n), err
}

// WriteFile creates or overwrites the file at the given path with the diagram.
func (u *UseCaseDiagram) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("uml.WriteFile: %w", err)
	}
	defer func() { _ = f.Close() }()
	_, err = u.WriteTo(f)
	return err
}
