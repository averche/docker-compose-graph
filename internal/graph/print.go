package graph

import (
	"fmt"
	"io"
	"strings"

	"github.com/averche/docker-compose-graph/internal/compose"
)

// Print will print the given nodes as a dot-graph
func Print(w io.Writer, nodes []Node) {
	fmt.Fprintf(w, `digraph compose {`+"\n")
	fmt.Fprintf(w, `  graph [fontname = "arial"];`+"\n")
	fmt.Fprintf(w, `  node  [fontname = "arial"];`+"\n")
	fmt.Fprintf(w, `  edge  [fontname = "arial" color = %q];`+"\n", DarkGrey)

	// appended to the names of subgraph clusters
	for _, node := range nodes {
		printNode(w, node.name, node.category, false, false)
	}

	printLegend(w, nodes)

	for _, node := range nodes {
		printDependencies(w, node.name, node.serviceDependencies)
	}

	for _, node := range nodes {
		printVolumeMounts(w, node.name, node.volumeMountsVolume)
	}

	fmt.Fprintf(w, "}")
}

// printLegend prints a dot-graph subgraph with all the node types we encountered
func printLegend(w io.Writer, nodes []Node) {
	fmt.Fprintf(w, "  subgraph cluster_0 {\n")
	fmt.Fprintf(w, "      label = %q\n", "Legend")
	fmt.Fprintf(w, "      shape = %q\n", Box)
	fmt.Fprintf(w, "      style = %q\n", JoinStyles([]Style{Rounded, Bold, Dashed}, ","))
	fmt.Fprintf(w, "      color = %q\n", DarkGrey)

	// ordered list of categories to achieve a reproducible output
	for _, category := range orderedPresentCategories(nodes) {
		printNode(w, category.String(), category, true, true)
	}

	fmt.Fprintf(w, "  }\n")
}

// printNode prints a dot-graph formatted string in the form 'name [style decorators];'
func printNode(w io.Writer, name string, category Category, offset, small bool) {
	d, ok := categoryDecorations[category]
	if !ok {
		panic(fmt.Sprintf("decorations missing for '%s' category", category))
	}

	var format string
	if small {
		format = `[shape = %-12q style = %-24q fillcolor = %-12q color = %-12q fontcolor = %-12q fontsize = "8pt"  label = %q];`
	} else {
		format = `[shape = %-12q style = %-24q fillcolor = %-12q color = %-12q fontcolor = %-12q label = %q];`
	}

	if offset {
		format = "    %-24s " + format + "\n"
	} else {
		format = "  %-26s " + format + "\n"
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
func printDependencies(w io.Writer, name string, dependencies []compose.ServiceDependency) {
	for _, dependency := range dependencies {
		switch dependency.Condition {
		case compose.ConditionServiceHealthy:
			fmt.Fprintf(w, `  %-26s -> %-26s [arrowhead="diamond" style="bold"];`+"\n", sanitize(name), sanitize(dependency.On))
		case compose.ConditionServiceCompletedSuccessfully:
			fmt.Fprintf(w, `  %-26s -> %-26s [style="bold"];`+"\n", sanitize(name), sanitize(dependency.On))
		default:
			fmt.Fprintf(w, `  %-26s -> %-26s [style="dashed"];`+"\n", sanitize(name), sanitize(dependency.On))
		}
	}
}

// printVolumeMounts prints the dependency lines formatted in dot-graph arrow (->) notation
func printVolumeMounts(w io.Writer, name string, volumeMounts []compose.VolumeMount) {
	for _, v := range volumeMounts {
		fmt.Fprintf(w, `  %-26s -> %-26s [arrowhead="diamond" style="bold"];`+"\n", sanitize(name), sanitize(v.Source))
	}
}

// dashes are not permitted in dot-graph names
func sanitize(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}
