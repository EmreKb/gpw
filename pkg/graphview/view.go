package graphview

import (
	"bytes"
	"context"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"

	"github.com/EmreKb/project-view/pkg/analyzer"
)

func View(pkg *analyzer.Package, outputPath string) error {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		panic(err)
	}

	graph, err := g.Graph(graphviz.WithName(pkg.Path))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()

	nodes := make(map[string]*cgraph.Node)

	rootNode, err := graph.CreateNodeByName(pkg.Path)
	if err != nil {
		return err
	}
	nodes[pkg.Path] = rootNode

	if err := mapNodes(graph, pkg, &nodes); err != nil {
		return err
	}

	if err := mapEdges(graph, pkg, &nodes); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		return err
	}

	if err := g.RenderFilename(ctx, graph, graphviz.PNG, outputPath); err != nil {
		panic(err)
	}

	return nil
}

func mapNodes(graph *cgraph.Graph, pkg *analyzer.Package, nodes *map[string]*cgraph.Node) error {
	for _, p := range pkg.Packages {
		node, err := graph.CreateNodeByName(p.Path)
		if err != nil {
			return err
		}

		(*nodes)[p.Path] = node

		mapNodes(graph, p, nodes)
	}

	return nil
}

func mapEdges(graph *cgraph.Graph, pkg *analyzer.Package, nodes *map[string]*cgraph.Node) error {
	for _, p := range pkg.Packages {
		_, err := graph.CreateEdgeByName(p.Path, (*nodes)[pkg.Path], (*nodes)[p.Path])
		if err != nil {
			return err
		}

		mapEdges(graph, p, nodes)
	}

	return nil
}
