package models

// Categories
type categories []struct {
	Name          string   `json:"name"`
	Subcategories []string `json:"subcategories"`
}

// Categories is where the music genres are defined
// We can use this to validate user input before performing database operations
// Since there are so few categories, there is no need to have this at the database layer
var Categories = categories{
	{
		Name:          "weed",
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
	for _, officialCategory := range Categories {
		if category == officialCategory.Name {
			for _, officialSubCategory := range officialCategory.Subcategories {
				if subcategory == officialSubCategory {
					return true
				}
			}
		}
	}
	return false
}
