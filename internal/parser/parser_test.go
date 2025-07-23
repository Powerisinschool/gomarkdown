package parser_test

import (
	"github.com/Powerisinschool/gomarkdown/internal/ast"
	"github.com/Powerisinschool/gomarkdown/internal/parser"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("Parse", func(t *testing.T) {
		testInput := "Hello, world!"
		p := parser.New([]byte(testInput))
		doc := p.Parse()
		if doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Text).Value != testInput {
			t.Errorf("expected '%s', got '%s'", testInput, doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Text).Value)
		}
	})
}

func TestParseEmpty(t *testing.T) {
	t.Run("ParseEmpty", func(t *testing.T) {
		testInput := ""
		p := parser.New([]byte(testInput))
		doc := p.Parse()
		if doc.Children != nil {
			t.Errorf("expected nil, got '%s'", doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Text).Value)
		}
	})
}

func TestParseHeading(t *testing.T) {
	t.Run("ParseHeading", func(t *testing.T) {
		testInput := "# Hello, world!"
		p := parser.New([]byte(testInput))
		doc := p.Parse()
		if doc.Children[0].(*ast.Heading).Children[0].(*ast.Text).Value != testInput[2:] {
			t.Errorf("expected '%s', got '%s'", testInput[2:], doc.Children[0].(*ast.Heading).Children[0].(*ast.Text).Value)
		}
	})
}

func TestParseFencedCodeBlock(t *testing.T) {
	t.Run("ParseFencedCodeBlock", func(t *testing.T) {
		testInput := "```bash  \nthis is some code\n\nhi\n```"
		p := parser.New([]byte(testInput))
		doc := p.Parse()
		codeBlock := doc.Children[0].(*ast.CodeBlock)
		if codeBlock.Language != "bash" {
			t.Errorf("expected 'bash', got '%s'", codeBlock.Language)
		}
		if codeBlock.Content != "this is some code\n\nhi\n" {
			t.Errorf("expected 'this is some code\n\nhi\n', got '%s'", codeBlock.Content)
		}
	})
}

func TestParseIndentedCodeBlock(t *testing.T) {
	t.Run("ParseIndentedCodeBlock", func(t *testing.T) {
		testInput := "\tthis is some code\n\t\n\thi\n"
		p := parser.New([]byte(testInput))
		doc := p.Parse()
		codeBlock := doc.Children[0].(*ast.CodeBlock)
		if codeBlock.Language != "" {
			t.Errorf("expected '', but got '%s'", codeBlock.Language)
		}
		if codeBlock.Content != "this is some code\n\nhi\n" {
			t.Errorf("expected 'this is some code\n\nhi\n', got '%s'", codeBlock.Content)
		}
	})
}

//func TestParseStrong(t *testing.T) {
//	t.Run("ParseStrong", func(t *testing.T) {
//		testInput := "**Hello, world!**"
//		p := parser.New([]byte(testInput))
//		doc := p.Parse()
//		if doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Strong).Children[0].(*ast.Text).Value != testInput {
//			t.Errorf("expected '%s', got '%s'", testInput, doc.Children[0].(*ast.Paragraph).Children[0].(*ast.Strong).Children[0].(*ast.Text).Value)
//		}
//	})
//}
