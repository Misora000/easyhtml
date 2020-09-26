package easyhtml

import (
	"bytes"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Tokenizer wraps golang.org/x/net/html with some useful methods.
type Tokenizer struct {
	z *html.Tokenizer
}

// Attrs is the attributes of a token.
type Attrs map[string]string

// NewTokenizer new a tokenizer.
func NewTokenizer(body io.Reader) *Tokenizer {
	z := html.NewTokenizer(body)
	z.NextIsNotRawText()
	return &Tokenizer{
		z: z,
	}
}

// JumpToTag move the tokenizer pointer to the target tag and return its attr.
func (z *Tokenizer) JumpToTag(tag string) (a Attrs, eof bool) {
	for {
		token := z.z.Next()
		if token == html.ErrorToken {
			eof = true
			return
		}

		name, _ := z.z.TagName()
		if string(name) != tag || token == html.EndTagToken {
			continue
		}

		a = z.getAttrs()
		break
	}
	return
}

// JumpToID move the tokenizer pointer to the target id and return its attr.
func (z *Tokenizer) JumpToID(tag, id string) (a Attrs, eof bool) {
	for {
		token := z.z.Next()
		if token == html.ErrorToken {
			eof = true
			return
		}

		name, hasAttr := z.z.TagName()
		if string(name) != tag || !hasAttr {
			continue
		}

		a = z.getAttrs()
		if val, exists := a["id"]; exists && val == id {
			break
		}
	}
	return
}

// JumpToClass move the tokenizer pointer to the target class and return its attr.
func (z *Tokenizer) JumpToClass(tag, class string) (a Attrs, eof bool) {
	for {
		token := z.z.Next()
		if token == html.ErrorToken {
			eof = true
			return
		}

		name, hasAttr := z.z.TagName()
		if string(name) != tag || !hasAttr {
			continue
		}

		a = z.getAttrs()
		val, exists := a["class"]
		if !exists {
			continue
		}

		for _, c := range strings.Split(val, " ") {
			if c == class {
				return
			}
		}
	}
}

// ExpandToken expand the current token as the raw text.
// It is useful for debug or starting a new tokenizer of this sub-DOM and
// doesn't affect the original tokenizer because it makes a copy form the
// original.
func (z *Tokenizer) ExpandToken() (*bytes.Buffer, bool) {
	buffer := new(bytes.Buffer)
	depth := 1

	for {
		switch z.z.Next() {
		case html.ErrorToken:
			return buffer, true

		case html.StartTagToken:
			name, _ := z.z.TagName()
			if computeDepth(string(name)) {
				depth++
			}

		case html.EndTagToken:
			name, _ := z.z.TagName()
			if computeDepth(string(name)) {
				depth--
			}
			if depth == 0 {
				return buffer, false
			}
		}
		buffer.Write(z.z.Raw())
	}
}

func computeDepth(name string) bool {
	// When computing depth, we only care about the following tags.
	allow := []string{"div", "span", "a", "li"}
	for _, a := range allow {
		if name == a {
			return true
		}
	}
	return false
}

func (z *Tokenizer) getAttrs() (o Attrs) {
	o = Attrs{}
	for {
		key, val, more := z.z.TagAttr()
		o[string(key)] = string(val)
		if !more {
			break
		}
	}
	return
}
