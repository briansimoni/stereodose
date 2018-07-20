package models

import "testing"

func Test_categories_Valid(t *testing.T) {
	type args struct {
		category    string
		subcategory string
	}
	tests := []struct {
		name string
		c    categories
		args args
		want bool
	}{
		{name: "Category does exist", want: true, c: Categories, args: args{category: "Weed", subcategory: "Thug Life"}},
		{name: "Category does not exist", want: false, c: Categories, args: args{category: "Weed", subcategory: "doesnotexist"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Valid(tt.args.category, tt.args.subcategory); got != tt.want {
				t.Errorf("categories.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
