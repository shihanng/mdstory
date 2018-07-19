// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/russross/blackfriday"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: publishRun,
	Args: cobra.ExactArgs(1),
}

var output string

func init() {
	rootCmd.AddCommand(publishCmd)

	publishCmd.Flags().StringVarP(&output, "output", "o", "", "Write the resulting HTML to file instead of stdout. If use with -xxx, the resulting HTML will not be posted to Medium.com but will be written to file.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func publishRun(cmd *cobra.Command, args []string) error {
	filename := args[0]

	blob, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "in reading markdown source")
	}

	var p *process

	if output != "" {
		var err error
		outputFile, err := os.Create(output)
		if err != nil {
			return errors.Wrap(err, "in creating output file")
		}
		defer outputFile.Close()

		p = &process{outputFile}
	} else {
		p = &process{os.Stdout}
	}

	p.parseMarkdown(blob)
	return nil
}

type process struct {
	output io.Writer
}

func (p *process) parseMarkdown(blob []byte) {
	md := blackfriday.New(blackfriday.WithExtensions(
		blackfriday.FencedCode | blackfriday.NoEmptyLineBeforeBlock))

	node := md.Parse(blob)

	w := &walker{}
	node.Walk(w.visitor)

	for _, image := range w.Images {
		rewriteImage(image)
	}

	for _, code := range w.CodeBlocks {
		replaceGist(code)
	}

	r := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{})

	node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return r.RenderNode(p.output, node, entering)
	})
}

type walker struct {
	CodeBlocks []*blackfriday.Node
	Images     []*blackfriday.Node
}

func (w *walker) visitor(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Image:
		w.Images = append(w.Images, node)
	case blackfriday.CodeBlock:
		w.CodeBlocks = append(w.CodeBlocks, node)
	}
	return blackfriday.GoToNext
}

func rewriteImage(imageNode *blackfriday.Node) {
	imageNode.LinkData.Destination = []byte(`test.jpg`)
}

func replaceGist(codeNode *blackfriday.Node) {
	n := blackfriday.NewNode(blackfriday.Link)
	n.LinkData.Destination = []byte(`example.com`)

	codeNode.InsertBefore(n)
}
