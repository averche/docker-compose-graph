package graph

import (
	"fmt"
	"io"
	"sort"
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
	var subgraphIndex uint32

	for _, node := range nodes {
		if len(node.volumes) != 0 {
			printNodeWithVolumes(w, node.name, node.category, subgraphIndex, node.volumes)
			subgraphIndex++
		} else {
			printNode(w, node.name, node.category, false, false)
		}
	}

	printLegend(w, nodes, subgraphIndex)

	for _, node := range nodes {
		printDependencies(w, node.name, node.dependencies)
	}

	fmt.Fprintf(w, "}")
}

// printLegend prints a dot-graph subgraph with all the node types we encountered
func printLegend(w io.Writer, nodes []Node, subgraphIndex uint32) {
	fmt.Fprintf(w, "  subgraph cluster_%d {\n", subgraphIndex)
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

// printNodeWithVolumes prints a dot-graph subgraph which includes the node and its volumes
func printNodeWithVolumes(w io.Writer, name string, category Category, subgraphIndex uint32, volumes []compose.Volume) {
	d, ok := categoryDecorations[category]
	if !ok {
		panic(fmt.Sprintf("decorations missing for '%s' category", category))
	}

	fmt.Fprintf(w, "  subgraph cluster_%d {\n", subgraphIndex)
	fmt.Fprintf(w, "      shape = %q\n", Box)
	fmt.Fprintf(w, "      style = %q\n", JoinStyles([]Style{Rounded, Bold, Dashed}, ","))
	fmt.Fprintf(w, "      color = %q\n", d.palette.ColorBorder)

	printNode(w, name, category, true, false)

	// sort the volumes to achieve a reproducible output
	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Source < volumes[j].Target
	})

	for i, volume := range volumes {
		fmt.Fprintf(
			w,
			`    %-24s [shape = %-12q style = %-24q                          color = %-12q fontcolor = %-12q label = "volume\nfrom: %s\nto: %s"];`+"\n",
			fmt.Sprintf("%s_v%d", sanitize(name), i),
			Cylinder,
			JoinStyles([]Style{Rounded, Bold, Dashed}, ","),
			d.palette.ColorBorder,
			DarkGrey,
			volume.Source,
			volume.Target,
		)
	}

	// draw connections to each volume
	for i := range volumes {
		fmt.Fprintf(w, `    %-24s -> %s_v%d;`+"\n", sanitize(name), sanitize(name), i)
	}

	fmt.Fprintf(w, "  }\n")
}

// printDependencies prints the dependency lines formatted in dot-graph arrow (->) notation
func printDependencies(w io.Writer, name string, dependencies []compose.Dependency) {
	// sort dependencies to achieve a reproducible output
	sort.Slice(dependencies, func(i, j int) bool {
		return dependencies[i].On < dependencies[j].On
	})

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

// dashes are not permitted in dot-graph names
func sanitize(name string) string {
	return strings.ReplaceAll(name, "-", "_")
}
