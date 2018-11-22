package models

// Categories
type categories map[string][]string

// Categories is where the music genres are defined
// We can use this to validate user input before performing database operations
// ... we could probably do this step on the database layer, but it works fine here
var Categories = categories{
	"weed":    []string{"chill", "groovin", "thug life"},
	"ecstacy": []string{"dance", "floored", "rolling balls"},
	"shrooms": []string{"matrix", "shaman", "space"},
	"lsd":     []string{"calm", "trippy", "fubar", "rockstar"},
}

func (c categories) Valid(category, subcategory string) bool {
	for _, sub := range c[category] {
		if sub == subcategory {
			return true
		}
	}
	return false
}
