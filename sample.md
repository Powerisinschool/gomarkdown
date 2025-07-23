# gomarkdown

A simple Markdown parser written in Go, built from scratch.

## Project Structure

- `cmd/gomarkdown`: The main entry point for the command-line tool.
- `internal/ast`: Defines the Abstract Syntax Tree (AST) for representing the markdown document structure.
- `internal/parser`: Contains the logic for parsing markdown text into an AST.
- `internal/renderer`: Handles the rendering of the AST into an output format, such as HTML.

## Getting Started

*This section will be updated as the project progresses.*

To build and run the project:

```bash
# (Instructions to be added)
```

[//]: # (**Note: *This *section* will be updated as the project progresses.* jkh**)

| Syntax      | Description | Test Text     |
| :---        |    :----:   |          ---: |
| Header      | Title       | Here's this   |
| Paragraph   | Text        | And more      |

|  Syntax   | Description |
|:---------:|-------------|
|  Header   | Title       |
| Paragraph | Text        |

&#124;

```json
{
  "firstName": "John",
  "lastName": "Smith",
  "age": 25
}
```

### My Great Heading {#custom-id}

Here's a simple footnote,[^1] and here's a longer one.[^bignote]

[^1]: This is the first footnote.

[^bignote]: Here's one with multiple paragraphs and code.

    Indent paragraphs to include them in the footnote.

    `{ my code }`

    Add as many paragraphs as you like.

[Heading IDs](#custom-id)

Gone camping! :tent: Be back soon.

That is so funny! :joy:

**In a few words:**

AI finds patterns in data to make predictions.

---

**In one sentence:**

AI systems learn from massive amounts
of information to recognize patterns
and then use them to make intelligent
decisions or predictions on their own.

ghdhgfdh **As a simple analogy:** gfhdhdfg

It's like teaching a child by showing it thousands of pictures of cats; eventually, it learns to recognize a cat it's never seen before.

*****djksbfjksbf*****
