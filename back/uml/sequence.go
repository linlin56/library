// Package uml provides a fluent builder to generate PlantUML
// sequence diagrams.
package uml

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// SequenceDiagram builds a PlantUML sequence diagram step by step.
type SequenceDiagram struct {
	title string
	lines []string
}

// NewSequenceDiagram creates a new diagram with an optional title.
func NewSequenceDiagram(title string) *SequenceDiagram {
	return &SequenceDiagram{title: title}
}

func (s *SequenceDiagram) add(l string) *SequenceDiagram {
	s.lines = append(s.lines, l)
	return s
}

// Participant declares a participant with an alias and a displayed label.
func (s *SequenceDiagram) Participant(alias, label string) *SequenceDiagram {
	return s.add(fmt.Sprintf(`participant "%s" as %s`, label, alias))
}

// Actor declares an actor (human representation).
func (s *SequenceDiagram) Actor(alias, label string) *SequenceDiagram {
	return s.add(fmt.Sprintf(`actor "%s" as %s`, label, alias))
}

// Arrow adds a solid arrow from from to to with a message.
func (s *SequenceDiagram) Arrow(from, to, msg string) *SequenceDiagram {
	return s.add(fmt.Sprintf("%s -> %s : %s", from, to, msg))
}

// DashedArrow adds a dashed arrow (response / return).
func (s *SequenceDiagram) DashedArrow(from, to, msg string) *SequenceDiagram {
	return s.add(fmt.Sprintf("%s --> %s : %s", from, to, msg))
}

// Activate marks a participant as active (lifeline bar).
func (s *SequenceDiagram) Activate(alias string) *SequenceDiagram {
	return s.add("activate " + alias)
}

// Deactivate removes a participant's lifeline bar.
func (s *SequenceDiagram) Deactivate(alias string) *SequenceDiagram {
	return s.add("deactivate " + alias)
}

// AltStart opens an alt/else/end block.
func (s *SequenceDiagram) AltStart(condition string) *SequenceDiagram {
	return s.add("alt " + condition)
}

// Else adds an alt branch.
func (s *SequenceDiagram) Else(condition string) *SequenceDiagram {
	return s.add("else " + condition)
}

// End closes an alt, group or loop block.
func (s *SequenceDiagram) End() *SequenceDiagram {
	return s.add("end")
}

// GroupStart opens a group block with a label.
func (s *SequenceDiagram) GroupStart(label string) *SequenceDiagram {
	return s.add("group " + label)
}

// Note adds a note next to a participant (side = "right of", "left of", "over").
func (s *SequenceDiagram) Note(side, alias, msg string) *SequenceDiagram {
	return s.add(fmt.Sprintf("note %s %s : %s", side, alias, msg))
}

// Separator adds a visual separator with a label.
func (s *SequenceDiagram) Separator(label string) *SequenceDiagram {
	return s.add("== " + label + " ==")
}

// Blank inserts an empty line for readability.
func (s *SequenceDiagram) Blank() *SequenceDiagram {
	return s.add("")
}

// String returns the full diagram in PlantUML syntax.
func (s *SequenceDiagram) String() string {
	var b strings.Builder
	b.WriteString("@startuml\n")
	if s.title != "" {
		fmt.Fprintf(&b, "title %s\n\n", s.title)
	}
	for _, l := range s.lines {
		b.WriteString(l + "\n")
	}
	b.WriteString("@enduml\n")
	return b.String()
}

// WriteTo writes the diagram to an io.Writer.
func (s *SequenceDiagram) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, s.String())
	return int64(n), err
}

// WriteFile creates (or overwrites) the file at the given path with the diagram.
func (s *SequenceDiagram) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("uml.WriteFile: %w", err)
	}
	defer func() { _ = f.Close() }()
	_, err = s.WriteTo(f)
	return err
}
