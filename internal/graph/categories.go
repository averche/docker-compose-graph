package graph

import (
	"regexp"
	"slices"
)

type Category uint8

const (
	CategoryNone Category = iota
	CategoryService
	CategoryVault
	CategoryCadence
	CategoryUserInterface
	CategoryTool
	CategoryDatabase
	CategoryStorage
	CategoryScript

	// this category is reserved for volumes
	CategoryVolume

	// this entry must be last
	categoryCount
)

var categoryStrings = []string{
	"none",
	"service",
	"vault",
	"cadence",
	"ui",
	"tool",
	"database",
	"storage",
	"script",
	"volume",
}

var categoryDecorations = map[Category]Decorations{
	CategoryNone:          {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Blue, DarkBlue, White}},
	CategoryService:       {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Blue, DarkBlue, White}},
	CategoryVault:         {styles: []Style{Rounded, Bold, Filled}, shape: Octagon, palette: Palette{Teal, DarkTeal, White}},
	CategoryCadence:       {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Red, DarkRed, White}},
	CategoryUserInterface: {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Purple, DarkPurple, White}},
	CategoryTool:          {styles: []Style{Rounded, Bold, Filled}, shape: Octagon, palette: Palette{Green, DarkGreen, White}},
	CategoryDatabase:      {styles: []Style{Rounded, Bold, Filled}, shape: Cylinder, palette: Palette{Green, DarkGreen, White}},
	CategoryStorage:       {styles: []Style{Rounded, Bold, Filled}, shape: Cylinder, palette: Palette{Red, DarkRed, White}},
	CategoryScript:        {styles: []Style{Bold, Filled}, shape: Note, palette: Palette{Grey, DarkGrey, White}},
	CategoryVolume:        {styles: []Style{Rounded, Bold, Filled}, shape: Cylinder, palette: Palette{Grey, DarkGrey, White}},
}

func (d Category) String() string {
	return categoryStrings[d]
}

// patterns are evaluated sequentially
var patterns = []struct {
	category Category
	pattern  *regexp.Regexp
}{{
	category: CategoryScript,
	pattern:  regexp.MustCompile(`(?i)^.*(script)$`),
}, {
	category: CategoryTool,
	pattern:  regexp.MustCompile(`(?i)^.*(tool)$`),
}, {
	category: CategoryStorage,
	pattern:  regexp.MustCompile(`(?i)^.*(storage)`),
}, {
	category: CategoryDatabase,
	pattern:  regexp.MustCompile(`(?i)^.*(database|postgres)`),
}, {
	category: CategoryUserInterface,
	pattern:  regexp.MustCompile(`(?i)^.*(ui)`),
}, {
	category: CategoryCadence,
	pattern:  regexp.MustCompile(`(?i)^.*(cadence|temporal)`),
}, {
	category: CategoryVault,
	pattern:  regexp.MustCompile(`(?i)^.*(vault)`),
}}

// DetermineCategory tries to guess the category based on the service name & label (if provided)
func DeterminteServiceCategory(service, label string) Category {
	// try to find exact matches for the label first
	if label != "" {
		if i := slices.Index(categoryStrings, label); i != -1 {
			return Category(i)
		}
	}

	// test for each category in sequence
	for _, p := range patterns {
		if p.pattern.MatchString(service) {
			return p.category
		}
	}

	return CategoryService // default
}
