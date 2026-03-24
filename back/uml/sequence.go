// Package uml fournit un builder fluide pour générer des diagrammes de
// séquence au format PlantUML.
package uml

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// SequenceDiagram construit un diagramme de séquence PlantUML pas à pas.
type SequenceDiagram struct {
	title string
	lines []string
}

// NewSequenceDiagram crée un nouveau diagramme avec un titre optionnel.
func NewSequenceDiagram(title string) *SequenceDiagram {
	return &SequenceDiagram{title: title}
}

func (s *SequenceDiagram) add(l string) *SequenceDiagram {
	s.lines = append(s.lines, l)
	return s
}

// Participant déclare un participant avec un alias et un libellé affiché.
func (s *SequenceDiagram) Participant(alias, label string) *SequenceDiagram {
	return s.add(fmt.Sprintf(`participant "%s" as %s`, label, alias))
}

// Actor déclare un acteur (représentation humaine).
func (s *SequenceDiagram) Actor(alias, label string) *SequenceDiagram {
	return s.add(fmt.Sprintf(`actor "%s" as %s`, label, alias))
}

// Arrow ajoute une flèche pleine de from vers to avec un message.
func (s *SequenceDiagram) Arrow(from, to, msg string) *SequenceDiagram {
	return s.add(fmt.Sprintf("%s -> %s : %s", from, to, msg))
}

// DashedArrow ajoute une flèche en pointillés (réponse / retour).
func (s *SequenceDiagram) DashedArrow(from, to, msg string) *SequenceDiagram {
	return s.add(fmt.Sprintf("%s --> %s : %s", from, to, msg))
}

// Activate marque un participant comme actif (barre de vie).
func (s *SequenceDiagram) Activate(alias string) *SequenceDiagram {
	return s.add("activate " + alias)
}

// Deactivate retire la barre de vie d'un participant.
func (s *SequenceDiagram) Deactivate(alias string) *SequenceDiagram {
	return s.add("deactivate " + alias)
}

// AltStart ouvre un bloc alt/else/end.
func (s *SequenceDiagram) AltStart(condition string) *SequenceDiagram {
	return s.add("alt " + condition)
}

// Else ajoute une branche alt.
func (s *SequenceDiagram) Else(condition string) *SequenceDiagram {
	return s.add("else " + condition)
}

// End ferme un bloc alt, group ou loop.
func (s *SequenceDiagram) End() *SequenceDiagram {
	return s.add("end")
}

// GroupStart ouvre un bloc group avec un libellé.
func (s *SequenceDiagram) GroupStart(label string) *SequenceDiagram {
	return s.add("group " + label)
}

// Note ajoute une note côté d'un participant (side = "right of", "left of", "over").
func (s *SequenceDiagram) Note(side, alias, msg string) *SequenceDiagram {
	return s.add(fmt.Sprintf("note %s %s : %s", side, alias, msg))
}

// Separator ajoute un séparateur visuel avec un libellé.
func (s *SequenceDiagram) Separator(label string) *SequenceDiagram {
	return s.add("== " + label + " ==")
}

// Blank insère une ligne vide pour l'aération.
func (s *SequenceDiagram) Blank() *SequenceDiagram {
	return s.add("")
}

// String retourne le diagramme complet en syntaxe PlantUML.
func (s *SequenceDiagram) String() string {
	var b strings.Builder
	b.WriteString("@startuml\n")
	if s.title != "" {
		b.WriteString(fmt.Sprintf("title %s\n\n", s.title))
	}
	for _, l := range s.lines {
		b.WriteString(l + "\n")
	}
	b.WriteString("@enduml\n")
	return b.String()
}

// WriteTo écrit le diagramme dans un io.Writer.
func (s *SequenceDiagram) WriteTo(w io.Writer) error {
	_, err := fmt.Fprint(w, s.String())
	return err
}

// WriteFile crée (ou écrase) le fichier au chemin donné avec le diagramme.
func (s *SequenceDiagram) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("uml.WriteFile: %w", err)
	}
	defer f.Close()
	return s.WriteTo(f)
}
