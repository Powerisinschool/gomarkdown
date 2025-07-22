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
	doc.Children = p.parseBlocks()
	//textNode := &ast.Text{Value: string(p.input)}
	//paragraphNode := &ast.Paragraph{
	//	Children: []ast.InlineNode{textNode},
	//}
	//doc.Children = []ast.BlockNode{paragraphNode}

	return doc
}

//func (p *Parser) parseLinesBlock() []ast.BlockNode {
//	rawBlocks := strings.Split(string(p.input), "\n\n")
//	blocks := make([]ast.BlockNode, len(rawBlocks))
//	for i, block := range rawBlocks {
//		block = strings.TrimSpace(block)
//		if block == "" {
//			continue
//		}
//		// Parse headings
//		if block[0] == '#' {
//			blk := strings.SplitN(block, " ", 2)
//			level := strings.Count(blk[0], "#")
//			if level >= 1 && level <= 6 && level == len(blk[0]) {
//				blocks[i] = &ast.Heading{Level: level, Children: p.inlineParse(blk[1])}
//				continue
//			}
//		}
//		// Default to a paragraph
//		if block != "" {
//			blocks[i] = &ast.Paragraph{Children: p.inlineParse(block)}
//		}
//	}
//	return blocks
//}

// parseList takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.List object.
func (p *Parser) parseList(lines []string, prefix string, i *int) *ast.List {
	listNode := &ast.List{}
	for *i < len(lines) {
		line := lines[*i]

		// If the line doesn't start with the separator, the list is over
		if !strings.HasPrefix(line, prefix) {
			break
		}
		// Get the content of the list item
		content := strings.TrimPrefix(line, prefix)
		// Run the content through the inline parser
		inlineContent := p.inlineParse(content)
		// Create the list item and add it to the list
		listItem := &ast.ListItem{Children: inlineContent}
		listNode.Items = append(listNode.Items, listItem)
		// Move to the next line
		*i++
	}
	return listNode
}

func (p *Parser) parseBlocks() []ast.BlockNode {
	var blocks []ast.BlockNode
	lines := strings.Split(string(p.input), "\n")

	for i := 0; i < len(lines); {
		line := lines[i]
		if strings.HasPrefix(line, "    ") {
			// TODO: Parse the whole code block and update `i`
			//codeBlockNode := p.parseCodeBlock(lines, &i)
			//blocks = append(blocks, codeBlockNode)
			//continue
			i++
			continue
		}
		line = strings.TrimSpace(line)
		// Check for unordered list items
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
			// Parse the whole list and update `i`
			listNode := p.parseList(lines, line[:2], &i)
			blocks = append(blocks, listNode)
			continue
		}
		// Check for ordered list items
		if strings.HasPrefix(line, "1. ") {
			// TODO: Parse the whole list and update `i`
			i++
			continue
		}
		// Skip empty lines
		if line == "" {
			i++
			continue
		}
		// Default to a paragraph
		pNode := &ast.Paragraph{Children: p.inlineParse(line)}
		blocks = append(blocks, pNode)
		i++
	}
	return blocks
}

func (p *Parser) inlineParse(line string) []ast.InlineNode {
	return []ast.InlineNode{&ast.Text{Value: strings.TrimSpace(line)}}
}
