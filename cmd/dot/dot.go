package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/zcong1993/go-graphviz"
	"github.com/zcong1993/go-graphviz/cgraph"
	"golang.org/x/crypto/ssh/terminal"
)

type Option struct {
	Format     graphviz.Format `description:"specify output format ( currently supported: dot svg png jpg )" short:"T" default:"dot"`
	Layout     graphviz.Layout `description:"specify layout engine ( currently supported: circo dot fdp neato nop nop1 nop2 osage patchwork sfdp twopi )" short:"K"`
	OutputFile string          `description:"specify output file name" short:"o" required:"true"`
}

func readGraph(args []string) (*cgraph.Graph, error) {
	if len(args) == 0 {
		if terminal.IsTerminal(0) {
			return nil, errors.New("required dot file or stdin")
		}
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		return graphviz.ParseBytes(bytes)
	}
	dotFile := args[0]
	return graphviz.ParseFile(dotFile)
}

func _main(args []string, opt *Option) (e error) {
	graph, err := readGraph(args)
	if err != nil {
		return err
	}
	g := graphviz.New()
	defer func() {
		if err := graph.Close(); err != nil {
			e = err
		}
		if err := g.Close(); err != nil {
			e = err
		}
	}()
	if opt.Layout != "" {
		g.SetLayout(opt.Layout)
	}
	return g.RenderFilename(graph, opt.Format, opt.OutputFile)
}

func main() {
	var opt Option
	parser := flags.NewParser(&opt, flags.Default)
	args, err := parser.Parse()
	if err != nil {
		return
	}
	if err := _main(args, &opt); err != nil {
		fmt.Println(err)
	}
}
