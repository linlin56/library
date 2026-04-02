// Package uml — UseCaseDiagram construit un diagramme de cas d'utilisation
// PlantUML de façon fluide, de manière similaire à SequenceDiagram.
package uml

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Diagram est une interface commune à tous les types de diagrammes PlantUML.
// Elle permet de stocker des types différents (séquence, cas d'utilisation…)
// dans une même map et de les écrire via une boucle unique.
type Diagram interface {
	WriteFile(path string) error
}

// UseCaseDiagram construit un diagramme de cas d'utilisation PlantUML.
type UseCaseDiagram struct {
	title string
	lines []string
}

// NewUseCaseDiagram crée un nouveau diagramme de cas d'utilisation avec un titre.
func NewUseCaseDiagram(title string) *UseCaseDiagram {
	return &UseCaseDiagram{title: title}
}

func (u *UseCaseDiagram) add(l string) *UseCaseDiagram {
	u.lines = append(u.lines, l)
	return u
}

// Direction définit l'orientation du diagramme (ex : "left to right direction").
func (u *UseCaseDiagram) Direction(dir string) *UseCaseDiagram {
	return u.add(dir)
}

// Actor déclare un acteur avec un alias et un libellé affiché.
func (u *UseCaseDiagram) Actor(alias, label string) *UseCaseDiagram {
	return u.add(fmt.Sprintf(`actor "%s" as %s`, label, alias))
}

// ActorInherits déclare une relation d'héritage entre acteurs (enfant hérite de parent).
// Produit : parent <|-- enfant
func (u *UseCaseDiagram) ActorInherits(child, parent string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s <|-- %s", parent, child))
}

// UseCase déclare un cas d'utilisation avec un alias et un libellé.
func (u *UseCaseDiagram) UseCase(alias, label string) *UseCaseDiagram {
	return u.add(fmt.Sprintf(`usecase "%s" as %s`, label, alias))
}

// RectangleStart ouvre un bloc rectangle (frontière système).
func (u *UseCaseDiagram) RectangleStart(label string) *UseCaseDiagram {
	return u.add(fmt.Sprintf(`rectangle "%s" {`, label))
}

// RectangleEnd ferme un bloc rectangle.
func (u *UseCaseDiagram) RectangleEnd() *UseCaseDiagram {
	return u.add("}")
}

// Association relie un acteur à un cas d'utilisation.
func (u *UseCaseDiagram) Association(from, to string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s --> %s", from, to))
}

// Include relie deux cas d'utilisation avec la relation <<include>>.
func (u *UseCaseDiagram) Include(from, to string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s ..> %s : <<include>>", from, to))
}

// Extend relie deux cas d'utilisation avec la relation <<extend>>.
func (u *UseCaseDiagram) Extend(from, to string) *UseCaseDiagram {
	return u.add(fmt.Sprintf("%s ..> %s : <<extend>>", from, to))
}

// Blank insère une ligne vide pour l'aération.
func (u *UseCaseDiagram) Blank() *UseCaseDiagram {
	return u.add("")
}

// String retourne le diagramme complet en syntaxe PlantUML.
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

// WriteTo écrit le diagramme dans un io.Writer.
func (u *UseCaseDiagram) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprint(w, u.String())
	return int64(n), err
}

// WriteFile crée ou écrase le fichier au chemin donné avec le diagramme.
func (u *UseCaseDiagram) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("uml.WriteFile: %w", err)
	}
	defer func() { _ = f.Close() }()
	_, err = u.WriteTo(f)
	return err
}
