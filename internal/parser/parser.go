package parser

import (
	"github.com/Powerisinschool/gomarkdown/internal/ast"
	"strings"
)

// Parser holds the state of the parser
type Parser struct {
	input []byte
}

// New creates a new Parser
func New(input []byte) *Parser {
	return &Parser{input: input}
}

// Parse is the main entry point for parsing
func (p *Parser) Parse() *ast.Document {
	// For now, we'll do the simplest possible thing:
	// treat the entire input as a single paragraph.
	doc := ast.NewDocument()

	// TODO: We will want to parse the input into a tree of
	// ast.InlineNode and ast.BlockNode objects.
	// TODO: We will want to parse the input line by line, rather than
	// treating the entire input as a single paragraph.
	doc.Children = p.parseLinesBlock()
	//textNode := &ast.Text{Value: string(p.input)}
	//paragraphNode := &ast.Paragraph{
	//	Children: []ast.InlineNode{textNode},
	//}
	//doc.Children = []ast.BlockNode{paragraphNode}

	return doc
}

func (p *Parser) parseLinesBlock() []ast.BlockNode {
	rawBlocks := strings.Split(string(p.input), "\n\n")
	blocks := make([]ast.BlockNode, len(rawBlocks))
	for i, block := range rawBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		// Parse headings
		if block[0] == '#' {
			blk := strings.SplitN(block, " ", 2)
			level := strings.Count(blk[0], "#")
			if level >= 1 && level <= 6 && level == len(blk[0]) {
				blocks[i] = &ast.Heading{Level: level, Children: p.inlineParse(blk[1])}
				continue
			}
		}
		// Default to a paragraph
		if block != "" {
			blocks[i] = &ast.Paragraph{Children: p.inlineParse(block)}
		}
	}
	return blocks
}

func (p *Parser) inlineParse(line string) []ast.InlineNode {
	return []ast.InlineNode{&ast.Text{Value: line}}
}
