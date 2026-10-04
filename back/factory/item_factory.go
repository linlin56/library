package factory

import "library/models"

// ItemFactory is the abstract factory.
type ItemFactory interface {
	CreateItem() models.IItem
}
