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
	doc := ast.NewDocument()
	doc.Children = p.parseBlocks()
	return doc
}

func (p *Parser) parseBlocks() []ast.BlockNode {
	var blocks []ast.BlockNode
	lines := strings.Split(string(p.input), "\n")

	for i := 0; i < len(lines); {
		line := lines[i]
		// Check for indented code blocks
		if strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			// TODO: Parse the whole code block and update `i`
			prefix := "    "
			if strings.HasPrefix(line, "\t") {
				prefix = "\t"
			}
			codeBlockNode := p.parseIndentedCodeBlock(lines, prefix, &i)
			blocks = append(blocks, codeBlockNode)
			continue
		}
		// Remove leading and trailing whitespace from the line
		line = strings.TrimSpace(line)
		// Check for fenced code blocks
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			codeBlockNode := p.parseFencedCodeBlock(lines, line[:3], &i)
			blocks = append(blocks, codeBlockNode)
			continue
		}
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
		// Parse headings
		if strings.HasPrefix(line, "#") {
			blk := strings.SplitN(line, " ", 2)
			level := strings.Count(blk[0], "#")
			if level >= 1 && level <= 6 && level == len(blk[0]) {
				blocks = append(blocks, &ast.Heading{Level: level, Children: p.inlineParse(blk[1])})
				i++
				continue
			}
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

// parseFencedCodeBlock takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.CodeBlock object.
func (p *Parser) parseFencedCodeBlock(lines []string, fence string, i *int) *ast.CodeBlock {
	// Safeguard against invalid input
	if !strings.HasPrefix(lines[*i], fence) {
		return &ast.CodeBlock{}
	}
	language := strings.TrimSpace(strings.TrimPrefix(lines[*i], fence))
	*i++

	codeBlock := &ast.CodeBlock{Language: language}
	for *i < len(lines) {
		line := lines[*i]
		if strings.TrimSpace(line) == fence {
			break
		}
		codeBlock.Content += line + "\n"
		*i++
	}
	// Ignore the rest of the content on the closing fence
	// Kept here in case of further modification as a safeguard
	// against infinite loops
	*i++
	return codeBlock
}

// parseIndentedCodeBlock takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.CodeBlock object.
func (p *Parser) parseIndentedCodeBlock(lines []string, prefix string, i *int) *ast.CodeBlock {
	codeBlock := &ast.CodeBlock{}
	for *i < len(lines) {
		line := lines[*i]
		if !strings.HasPrefix(line, prefix) {
			break
		}
		codeBlock.Content += strings.TrimPrefix(line, prefix) + "\n"
		*i++
	}
	return codeBlock
}
