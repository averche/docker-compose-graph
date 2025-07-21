package graph

import "regexp"

type Category uint8

const (
	CategoryNone Category = iota
	CategoryVolume
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

var categoryStrings = [...]string{
	"none",
	"Volume",
	"Service",
	"Vault",
	"Cadence",
	"FrontEnd",
	"Proxy",
	"Database",
	"Storage",
	"Script",
}

var categoryDecorations = map[Category]Decorations{
	CategoryNone:     {styles: []Style{Rounded, Bold, Filled}, shape: Box, palette: Palette{Blue, DarkBlue, White}},
	CategoryVolume:   {styles: []Style{Rounded, Bold, Filled}, shape: Cylinder, palette: Palette{Grey, DarkGrey, White}},
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
