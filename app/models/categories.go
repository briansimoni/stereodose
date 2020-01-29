package models

// Categories
type categories map[string][]string
type categories2 []struct {
	Name          string   `json:"name"`
	Subcategories []string `json:"subcategories"`
}

// Categories is where the music genres are defined
// We can use this to validate user input before performing database operations
// Since there are so few categories, there is no need to have this at the database layer
var Categories = categories{
	"weed":    []string{"chill", "groovin", "thug life"},
	"ecstacy": []string{"dance", "floored", "rolling balls"},
	"shrooms": []string{"matrix", "shaman", "space"},
	"lsd":     []string{"calm", "trippy", "rockstar"},
}

var Categories2 = categories2{
	{
		Name:          "Weed",
		Subcategories: []string{"chill", "groovin", "thug life"},
	},
	{
		Name:          "ecstacy",
		Subcategories: []string{"dance", "floored", "rolling balls"},
	},
	{
		Name:          "shrooms",
		Subcategories: []string{"matrix", "shaman", "space"},
	},
	{
		Name:          "lsd",
		Subcategories: []string{"calm", "trippy", "rockstar"},
	},
}

func (c categories) Valid(category, subcategory string) bool {
	for _, sub := range c[category] {
		if sub == subcategory {
			return true
		}
	}
	return false
}
