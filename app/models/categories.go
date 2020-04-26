package models

import (
	"fmt"
)

// Categories
type categories []category

type category struct {
	DisplayName   string   `json:"displayName"`
	Name          string   `json:"name"`
	Subcategories []string `json:"subcategories"`
}

// Categories is where the music genres are defined
// We can use this to validate user input before performing database operations
// Since there are so few categories, there is no need to have this at the database layer
var Categories = categories{
	{
		DisplayName:   "Weed",
		Name:          "weed",
		Subcategories: []string{"chill", "groovin", "thug life"},
	},
	{
		DisplayName:   "Ecstacy",
		Name:          "ecstacy",
		Subcategories: []string{"dance", "floored", "rolling balls"},
	},
	{
		DisplayName:   "Shrooms",
		Name:          "shrooms",
		Subcategories: []string{"matrix", "shaman", "space"},
	},
	{
		DisplayName:   "LSD",
		Name:          "LSD",
		Subcategories: []string{"calm", "trippy", "rockstar"},
	},
}

func getCategoryFromName(categoryName string) (category, error) {
	for _, cat := range Categories {
		if cat.Name == categoryName {
			return cat, nil
		}
	}
	return category{}, fmt.Errorf("Unable to find category with name %s", categoryName)
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
