package renderer

import (
	"github.com/Powerisinschool/gomarkdown/internal/ast"

	"io"
	"strconv"
)

// Renderer holds the configuration for rendering
type Renderer struct {
	w io.Writer // The destination for the HTML output
}

func New(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

func (r *Renderer) Render(node ast.Node) error {
	return r.renderNode(node)
}

// renderNode is the core recursive function that walks the AST and renders it
func (r *Renderer) renderNode(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Document:
		for _, c := range n.Children {
			if err := r.renderNode(c); err != nil {
				return err
			}
		}
	case *ast.List:
		if _, err := r.w.Write([]byte("<ul>")); err != nil {
			return err
		}
		for _, c := range n.Items {
			if err := r.renderNode(c); err != nil {
				return err
			}
		}
		if _, err := r.w.Write([]byte("</ul>\n")); err != nil {
			return err
		}
	case *ast.ListItem:
		if _, err := r.w.Write([]byte("<li>")); err != nil {
			return err
		}
		for _, c := range n.Children {
			if err := r.renderNode(c); err != nil {
				return err
			}
		}
		if _, err := r.w.Write([]byte("</li>\n")); err != nil {
			return err
		}
	case *ast.Heading:
		if _, err := r.w.Write([]byte("<h" + strconv.Itoa(int(n.Level)) + ">")); err != nil {
			return err
		}
		for _, c := range n.Children {
			if err := r.renderNode(c); err != nil {
				return err
			}
		}
		if _, err := r.w.Write([]byte("</h" + strconv.Itoa(int(n.Level)) + ">\n")); err != nil {
			return err
		}
	case *ast.Paragraph:
		// TODO: We will want to trim trailing and leading space
		// from paragraphs.
		if _, err := r.w.Write([]byte("<p>")); err != nil {
			return err
		}
		for _, c := range n.Children {
			if err := r.renderNode(c); err != nil {
				return err
			}
		}
		if _, err := r.w.Write([]byte("</p>\n")); err != nil {
			return err
		}
	case *ast.Text:
		_, err := r.w.Write([]byte(n.Value))
		if err != nil {
			return err
		}
	}
	return nil
}
