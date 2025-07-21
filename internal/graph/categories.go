package graph

import "regexp"

type Category uint8

const (
	CategoryNone Category = iota
	CategoryService
	CategoryVault
	CategoryCadence
	CategoryFrontEnd
	CategoryProxy
	CategoryDatabase
	CategoryStorage
	CategoryScript

	// this entry must be last
	categoryCount
)

var decorations = map[Category]Decorations{
	CategoryService:  {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Blue, DarkBlue, White}},
	CategoryVault:    {styles: []Style{Rounded, Bold, Filled}, shape: Record, palette: Palette{Teal, DarkTeal, White}},
	CategoryCadence:  {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Teal, DarkTeal, White}},
	CategoryFrontEnd: {styles: []Style{Rounded, Bold, Filled}, shape: Record, palette: Palette{Teal, DarkTeal, White}},
	CategoryProxy:    {styles: []Style{Rounded, Bold, Filled}, shape: Diamond, palette: Palette{Purple, DarkPurple, White}},
	CategoryDatabase: {styles: []Style{Rounded, Bold, Filled}, shape: Cylinder, palette: Palette{Green, DarkGreen, White}},
	CategoryStorage:  {styles: []Style{Rounded, Bold, Filled}, shape: Cylinder, palette: Palette{Red, DarkRed, White}},
	CategoryScript:   {styles: []Style{Rounded, Bold, Filled}, shape: Hexagon, palette: Palette{Grey, DarkGrey, White}},
}

func (d Category) String() string {
	return [...]string{
		"none",
		"Service",
		"Vault",
		"Cadence",
		"FrontEnd",
		"Proxy",
		"Database",
		"Storage",
		"Script",
	}[d]
}

type entity struct {
	name     string
	category Category
}

// patterns are evaluated sequentially
var patterns = []struct {
	category Category
	pattern  *regexp.Regexp
}{{
	category: CategoryScript,
	pattern:  regexp.MustCompile(`(?i)^.*(script)$`),
}, {
	category: CategoryProxy,
	pattern:  regexp.MustCompile(`(?i)^.*(proxy)$`),
}, {
	category: CategoryStorage,
	pattern:  regexp.MustCompile(`(?i)^.*(s3|storage)`),
}, {
	category: CategoryDatabase,
	pattern:  regexp.MustCompile(`(?i)^.*(database|postgres)`),
}, {
	category: CategoryFrontEnd,
	pattern:  regexp.MustCompile(`(?i)^.*(ui)`),
}, {
	category: CategoryCadence,
	pattern:  regexp.MustCompile(`(?i)^.*(cadence|temporal)`),
}, {
	category: CategoryVault,
	pattern:  regexp.MustCompile(`(?i)^.*(vault)`),
}}

// DetermineCategory tries to guess the category of the given thing based on the regex expressions above
func DeterminteCategory(thing string) Category {
	// test for each category in sequence
	for _, p := range patterns {
		if p.pattern.MatchString(thing) {
			return p.category
		}
	}

	return CategoryService // default
}

// OrderedPresentCategories returns an ordered list of categories that are present in the given slice
func OrderedPresentCategories(entities []entity) []Category {
	// bitmap intexed by category
	var exists [categoryCount]bool

	for _, e := range entities {
		exists[int(e.category)] = true
	}

	var present []Category

	for category := CategoryNone; category < categoryCount; category++ {
		if exists[category] {
			present = append(present, category)
		}
	}

	return present
}
