package models

import (
	"reflect"
	"testing"
)

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
		{name: "Category does exist", want: true, c: Categories, args: args{category: "weed", subcategory: "thug life"}},
		{name: "Category does not exist", want: false, c: Categories, args: args{category: "weed", subcategory: "doesnotexist"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Valid(tt.args.category, tt.args.subcategory); got != tt.want {
				t.Errorf("categories.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCategoryFromName(t *testing.T) {
	type args struct {
		categoryName string
	}
	tests := []struct {
		name    string
		args    args
		want    category
		wantErr bool
	}{
		{
			name: "valid test",
			args: args{categoryName: "weed"},
			want: category{
				DisplayName:   "Weed",
				Name:          "weed",
				Subcategories: []string{"chill", "groovin", "thug life"},
			},
			wantErr: false,
		},
		{
			name:    "invalid test",
			args:    args{categoryName: "poop"},
			want:    category{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCategoryFromName(tt.args.categoryName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCategoryFromName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCategoryFromName() = %v, want %v", got, tt.want)
			}
		})
	}
}
