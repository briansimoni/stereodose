package models

// Categories
type categories map[string][]string

// Categories is where the music genres are defined
// We can use this to validate user input before performing database operations
// ... we could probably do this step on the database layer, but it works fine here
var Categories = categories{
	"Weed":    []string{"Chill", "Groovin", "Thug Life"},
	"Ecstacy": []string{"clouds", "unicorns", "rainbows"},
	"Shrooms": []string{"mario", "luigi", "wario"},
	"LSD":     []string{"trippy1", "trippy2", "trippy3"},
}

func (c categories) Valid(category, subcategory string) bool {
	for _, sub := range c[category] {
		if sub == subcategory {
			return true
		}
	}
	return false
}
