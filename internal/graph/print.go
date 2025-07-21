package graph

import (
	"fmt"
	"io"
	"strings"

	"github.com/averche/docker-compose-graph/internal/compose"
)

// Print will print the given nodes as a dot-graph
func Print(w io.Writer, groups []NodeGroup) {
	fmt.Fprintf(w, `digraph compose {`+"\n")
	fmt.Fprintf(w, `  graph [fontname = "arial"];`+"\n")
	fmt.Fprintf(w, `  node  [fontname = "arial"];`+"\n")
	fmt.Fprintf(w, `  edge  [fontname = "arial" color = %q];`+"\n", DarkGrey)

	// subgraphIndex is appended to the names of subgraph clusters
	var subgraphIndex uint32

	for _, group := range groups {
		printGroups(w, group, subgraphIndex)
		subgraphIndex++
	}

	printLegend(w, groups, subgraphIndex)

	for _, group := range groups {
		for _, node := range group.Nodes {
			printDependencies(w, node.Name, node.ServiceDependencies, node.VolumeMounts)
		}
	}

	fmt.Fprintf(w, "}")
}

func printGroups(w io.Writer, group NodeGroup, subgraphIndex uint32) {
	fmt.Fprintf(w, "  subgraph cluster_%d {\n", subgraphIndex)
	fmt.Fprintf(w, "      label = %q\n", group.Name)
	fmt.Fprintf(w, "      shape = %q\n", Box)
	fmt.Fprintf(w, "      style = %q\n", JoinStyles([]Style{Rounded, Bold, Dashed}, ","))
	fmt.Fprintf(w, "      color = %q\n", DarkGrey)

	for _, node := range group.Nodes {
		printNode(w, node.Name, node.Category, false)
	}

	fmt.Fprintf(w, "  }\n")
}

// printLegend prints a dot-graph subgraph with all the node types we encountered
func printLegend(w io.Writer, groups []NodeGroup, subgraphIndex uint32) {
	fmt.Fprintf(w, "  subgraph cluster_%d {\n", subgraphIndex)
	fmt.Fprintf(w, "      label = %q\n", "Legend")
	fmt.Fprintf(w, "      shape = %q\n", Box)
	fmt.Fprintf(w, "      style = %q\n", JoinStyles([]Style{Rounded, Bold, Dashed}, ","))
	fmt.Fprintf(w, "      color = %q\n", DarkGrey)

	// ordered list of categories to achieve a reproducible output
	for _, category := range orderedPresentCategories(groups) {
		printNode(w, category.String(), category, true)
	}

	fmt.Fprintf(w, "  }\n")
}

// printNode prints a dot-graph formatted string in the form 'name [style decorators];'
func printNode(w io.Writer, name string, category Category, small bool) {
	d, ok := categoryDecorations[category]
	if !ok {
		panic(fmt.Sprintf("decorations missing for '%s' category", category))
	}

	var format string
	if small {
		format = `    %-34s [shape = %-12q style = %-24q fillcolor = %-12q color = %-12q fontcolor = %-12q fontsize = "8pt"  label = %q];` + "\n"
	} else {
		format = `    %-34s [shape = %-12q style = %-24q fillcolor = %-12q color = %-12q fontcolor = %-12q label = %q];` + "\n"
	}

	fmt.Fprintf(
		w,
		format,
		sanitize(name),
		d.shape,
		JoinStyles(d.styles, ","),
		d.palette.ColorFill,
		d.palette.ColorBorder,
		d.palette.ColorFont,
		name,
	)
}

// printDependencies prints the dependency lines formatted in dot-graph arrow (->) notation
func printDependencies(w io.Writer, name string, serviceDependencies []compose.ServiceDependency, volumeMounts []compose.VolumeMount) {
	for _, dependency := range serviceDependencies {
		switch dependency.Condition {
		case compose.ConditionServiceHealthy:
			fmt.Fprintf(w, `  %-36s -> %-36s [arrowhead="diamond" style="bold"];`+"\n", sanitize(name), sanitize(dependency.On))

		case compose.ConditionServiceCompletedSuccessfully:
			fmt.Fprintf(w, `  %-36s -> %-36s [style="bold"];`+"\n", sanitize(name), sanitize(dependency.On))

		case compose.ConditionServiceStarted:
			fmt.Fprintf(w, `  %-36s -> %-36s [style="dashed"];`+"\n", sanitize(name), sanitize(dependency.On))

		default:
			panic(fmt.Sprintf("unexpected dependency condition %q", dependency.Condition))
		}
	}

	for _, v := range volumeMounts {
		if v.ReadOnly {
			fmt.Fprintf(w, `  %-36s -> %-36s [style="dashed"];`+"\n", sanitize(name), sanitize(v.Source))
		} else {
			fmt.Fprintf(w, `  %-36s -> %-36s [style="bold"];`+"\n", sanitize(name), sanitize(v.Source))
		}
	}
}

// dashes are not permitted in dot-graph names
func sanitize(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}
