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

// isBlockBoundary checks if a line starts a new block, terminating a paragraph.
func (p *Parser) isBlockBoundary(line string) bool {
	trimmedLine := strings.TrimSpace(line)
	// Empty line is a boundary
	if trimmedLine == "" {
		return true
	}
	// Check for prefixes of other block types
	return strings.HasPrefix(trimmedLine, "#") ||
		strings.HasPrefix(trimmedLine, "- ") ||
		strings.HasPrefix(trimmedLine, "* ") ||
		strings.HasPrefix(trimmedLine, "+ ") ||
		strings.HasPrefix(trimmedLine, "```") ||
		strings.HasPrefix(trimmedLine, "~~~") ||
		(len(trimmedLine) > 0 && trimmedLine[0] >= '0' && trimmedLine[0] <= '9' && strings.HasPrefix(trimmedLine[1:], ". "))
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
			blocks = append(blocks, p.parseFencedCodeBlock(lines, line[:3], &i))
			continue
		}
		// Check for unordered list items
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") || strings.HasPrefix(line, "+ ") {
			blocks = append(blocks, p.parseList(lines, &i))
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
		paragraphNode := p.parseParagraph(lines, &i)
		//pNode := &ast.Paragraph{Children: p.inlineParse(line)}
		blocks = append(blocks, paragraphNode)
	}
	return blocks
}

// inlineParse takes in a line of text and returns a slice of inline nodes
func (p *Parser) inlineParse(line string) []ast.InlineNode {
	// Pass 1: Tokenize the line into text and delimiters
	tokens := p.tokenizeInline(line)
	// Pass 2. Process the tokens to resolve delimiters into Strong/Emphasis nodes.
	return p.processTokens(tokens)
}

// parseList takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.List object.
func (p *Parser) parseList(lines []string, i *int) *ast.List {
	listNode := &ast.List{}
	// Continue parsing as long as the lines could be part of the current list.
	for *i < len(lines) {
		line := lines[*i]
		trimmedLine := strings.TrimSpace(line)
		prefix := ""

		if strings.HasPrefix(trimmedLine, "- ") {
			prefix = "- "
		} else if strings.HasPrefix(trimmedLine, "* ") {
			prefix = "* "
		} else if strings.HasPrefix(trimmedLine, "+ ") {
			prefix = "+ "
		} else {
			// Not a list item, so the list ends here.
			break
		}

		content := strings.TrimPrefix(trimmedLine, prefix)
		// Basic multi-line handling (a more advanced parser would handle indentation).
		// For now, we assume one line per list item for simplicity.
		listItem := &ast.ListItem{Children: p.inlineParse(content)}
		listNode.Items = append(listNode.Items, listItem)
		*i++
	}
	return listNode
}

// parseFencedCodeBlock takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.CodeBlock object.
func (p *Parser) parseFencedCodeBlock(lines []string, fence string, i *int) *ast.CodeBlock {
	firstLine := strings.TrimSpace(lines[*i])
	language := strings.TrimSpace(strings.TrimPrefix(firstLine, fence))
	*i++

	var content strings.Builder
	for *i < len(lines) {
		line := lines[*i]
		if strings.TrimSpace(line) == fence {
			*i++ // Consume the closing fence line
			break
		}
		if content.Len() > 0 {
			content.WriteString("\n")
		}
		content.WriteString(line)
		*i++
	}
	return &ast.CodeBlock{Language: language, Content: content.String()}
}

// parseIndentedCodeBlock takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.CodeBlock object.
func (p *Parser) parseIndentedCodeBlock(lines []string, prefix string, i *int) *ast.CodeBlock {
	var content strings.Builder
	for *i < len(lines) {
		line := lines[*i]
		if !strings.HasPrefix(line, prefix) && strings.TrimSpace(line) != "" {
			break
		}
		if content.Len() > 0 {
			content.WriteString("\n")
		}
		content.WriteString(strings.TrimPrefix(line, prefix))
		*i++
	}
	return &ast.CodeBlock{Content: content.String()}
}

// parseParagraph now correctly parses multi-line paragraphs.
func (p *Parser) parseParagraph(lines []string, i *int) *ast.Paragraph {
	var content strings.Builder
	// Consume lines until a block boundary is found.
	for *i < len(lines) {
		line := lines[*i]
		if p.isBlockBoundary(line) {
			break
		}
		if content.Len() > 0 {
			content.WriteString(" ")
		}
		content.WriteString(strings.TrimSpace(line))
		*i++
	}
	return &ast.Paragraph{Children: p.inlineParse(content.String())}
}

func (p *Parser) tokenizeInline(line string) []ast.InlineNode {
	var nodes []ast.InlineNode
	textBuffer := &strings.Builder{}
	i := 0
	for i < len(line) {
		// For this simplified version, we'll focus on '*' as the delimiter.
		// A full parser would also handle ` and other special characters here.
		if line[i] == '*' {
			if textBuffer.Len() > 0 {
				nodes = append(nodes, &ast.Text{Value: textBuffer.String()})
				textBuffer.Reset()
			}
			start := i
			for i < len(line) && line[i] == '*' {
				i++
			}
			count := i - start
			nodes = append(nodes, &ast.Delimiter{Value: line[start:i], Count: count})
			continue // Continue to the next part of the string
		}
		textBuffer.WriteByte(line[i])
		i++
	}
	if textBuffer.Len() > 0 {
		nodes = append(nodes, &ast.Text{Value: textBuffer.String()})
	}
	return nodes
}

func (p *Parser) processTokens(tokens []ast.InlineNode) []ast.InlineNode {
	// A stack to keep track of the indices of opening delimiters.
	openerStack := []int{}

	// Iterate through tokens by index.
	for i := 0; i < len(tokens); i++ {
		closer, isCloser := tokens[i].(*ast.Delimiter)
		if !isCloser {
			continue // Skip non-delimiter tokens
		}

		// Search backwards through the stack for a matching opener.
		var openerIndexInStack int = -1
		var openerToken *ast.Delimiter
		for j := len(openerStack) - 1; j >= 0; j-- {
			// A closer can match an opener of the same type (`*` or `**`).
			// This simplified logic checks if both can form a strong or emphasis pair.
			openerCand := tokens[openerStack[j]].(*ast.Delimiter)
			if (openerCand.Count >= 2 && closer.Count >= 2) || (openerCand.Count >= 1 && closer.Count >= 1) {
				openerIndexInStack = j
				openerToken = openerCand
				break
			}
		}

		if openerToken != nil {
			// Found a valid pair.
			openerIndexInTokens := openerStack[openerIndexInStack]

			// Determine how many asterisks to use (the "stealing" logic).
			numToUse := 2
			if openerToken.Count < 2 || closer.Count < 2 {
				numToUse = 1
			}

			// Create the new node (Strong or Emphasis).
			var newNode ast.InlineNode
			// The children are the tokens between the opener and closer.
			children := p.processTokens(tokens[openerIndexInTokens+1 : i])
			if numToUse == 2 {
				newNode = &ast.Strong{Children: children}
			} else {
				newNode = &ast.Emphasis{Children: children}
			}

			// "Consume" the delimiters by reducing their counts.
			openerToken.Count -= numToUse
			closer.Count -= numToUse

			// Replace the entire sequence (opener, children, closer) with the new node.
			// First, create a new slice containing elements before the opener.
			newTokens := append([]ast.InlineNode{}, tokens[:openerIndexInTokens]...)
			// Append the new combined node.
			newTokens = append(newTokens, newNode)
			// Append the elements after the closer.
			newTokens = append(newTokens, tokens[i+1:]...)

			// Replace the old token slice and restart the scan.
			tokens = newTokens
			// Reset index to re-evaluate from the point of insertion.
			i = openerIndexInTokens

			// Clear the opener stack as the context has changed.
			// A more optimized solution could rebuild the stack, but this is safer.
			openerStack = []int{}
		} else {
			// No matching opener found, so this is a potential opener.
			openerStack = append(openerStack, i)
		}
	}

	// Final pass: convert any leftover delimiters back to text nodes.
	var finalNodes []ast.InlineNode
	for _, token := range tokens {
		if delim, ok := token.(*ast.Delimiter); ok && delim.Count > 0 {
			finalNodes = append(finalNodes, &ast.Text{Value: delim.Value})
		} else if _, ok := token.(*ast.Delimiter); !ok {
			finalNodes = append(finalNodes, token)
		}
	}

	return finalNodes
}

//func (p *Parser) buildTree(nodes []ast.InlineNode, matches map[*ast.Delimiter]*ast.Delimiter) []ast.InlineNode {
//	var result []ast.InlineNode
//	i := 0
//	for i < len(nodes) {
//		node := nodes[i]
//		opener, isOpener := node.(*ast.Delimiter)
//
//		if isOpener && matches[opener] != nil {
//			closer := matches[opener]
//			// Find the closer's index
//			closerIndex := -1
//			for j := i + 1; j < len(nodes); j++ {
//				if nodes[j] == closer {
//					closerIndex = j
//					break
//				}
//			}
//
//			// Determine if it's strong or emphasis
//			count := 1
//			if opener.Count >= 2 && closer.Count >= 2 {
//				count = 2
//			}
//
//			// Recursively build the children
//			children := p.buildTree(nodes[i+1:closerIndex], matches)
//			var newNode ast.InlineNode
//			if count == 1 {
//				newNode = &ast.Emphasis{Children: children}
//			} else {
//				newNode = &ast.Strong{Children: children}
//			}
//			result = append(result, newNode)
//
//			// Jump the index past the closer
//			i = closerIndex + 1
//		} else {
//			// Not an opener or not matched, so just add the node.
//			// If it's an unmatched delimiter, convert it to text.
//			if delim, ok := node.(*ast.Delimiter); ok {
//				result = append(result, &ast.Text{Value: delim.Value})
//			} else {
//				result = append(result, node)
//			}
//			i++
//		}
//	}
//	return result
//}

// parseList takes in a slice of lines and a pointer to the current index.
// It returns a pointer to an ast.List object.
//func (p *Parser) parseList(lines []string, prefix string, i *int) *ast.List {
//	listNode := &ast.List{}
//	for *i < len(lines) {
//		line := lines[*i]
//
//		// If the line doesn't start with the separator, the list is over
//		if !strings.HasPrefix(line, prefix) {
//			break
//		}
//		// Get the content of the list item
//		content := strings.TrimPrefix(line, prefix)
//		// Run the content through the inline parser
//		inlineContent := p.inlineParse(content)
//		// Create the list item and add it to the list
//		listItem := &ast.ListItem{Children: inlineContent}
//		listNode.Items = append(listNode.Items, listItem)
//		// Move to the next line
//		*i++
//	}
//	return listNode
//}
