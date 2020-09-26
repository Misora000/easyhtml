package easyhtml

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

const rawHTML = `<!DOCTYPE html>
<html>
<title>Title</title>
<body>
<div>
	<p>Main</p>
	<!-- pictures -->
	<img src="123.jpeg" class="small-pic" id="pic1" data-bind="hello" />
	<img src="456.jpeg" class="medium-pic" id="pic2" data-bind="world" />
	<img src="789.jpeg" class="small-pic" id="pic3" data-bind="TENET" />
	<a href="#test">Test</a>
</div>
<li>
	<ol>option 1</ol>
	<ol>option 2</ol>
	<ol>option 3</ol>
</li>
</body>
</html>`

func (s *testSuite) TestJumpToID() {
	z := NewTokenizer(strings.NewReader(rawHTML))

	attr, eof := z.JumpToID("img", "pic1")
	s.Require().False(eof)
	s.Require().Len(attr, 4)
	for i, v := range attr {
		switch i {
		case "src":
			s.Require().Equal("123.jpeg", v)
		case "class":
			s.Require().Equal("small-pic", v)
		case "data-bind":
			s.Require().Equal("hello", v)
		case "id":
			s.Require().Equal("pic1", v)
		default:
			s.Require().Fail("unexpected attr: %v=%v", i, v)
		}
	}

	attr, eof = z.JumpToID("img", "pic3")
	s.Require().False(eof)
	s.Require().Len(attr, 4)
	for i, v := range attr {
		switch i {
		case "src":
			s.Require().Equal("789.jpeg", v)
		case "class":
			s.Require().Equal("small-pic", v)
		case "data-bind":
			s.Require().Equal("TENET", v)
		case "id":
			s.Require().Equal("pic3", v)
		default:
			s.Require().Fail("unexpected attr: %v=%v", i, v)
		}
	}

	// Not found until EOF.
	_, eof = z.JumpToID("img", "pic3")
	s.Require().True(eof)
}

func (s *testSuite) TestJumpToClass() {
	z := NewTokenizer(strings.NewReader(rawHTML))

	attr, eof := z.JumpToClass("img", "small-pic")
	s.Require().False(eof)
	s.Require().Len(attr, 4)
	for i, v := range attr {
		switch i {
		case "src":
			s.Require().Equal("123.jpeg", v)
		case "class":
			s.Require().Equal("small-pic", v)
		case "data-bind":
			s.Require().Equal("hello", v)
		case "id":
			s.Require().Equal("pic1", v)
		default:
			s.Require().Fail("unexpected attr: %v=%v", i, v)
		}
	}

	attr, eof = z.JumpToClass("img", "small-pic")
	s.Require().False(eof)
	s.Require().Len(attr, 4)
	for i, v := range attr {
		switch i {
		case "src":
			s.Require().Equal("789.jpeg", v)
		case "class":
			s.Require().Equal("small-pic", v)
		case "data-bind":
			s.Require().Equal("TENET", v)
		case "id":
			s.Require().Equal("pic3", v)
		default:
			s.Require().Fail("unexpected attr: %v=%v", i, v)
		}
	}

	// Not found until EOF.
	_, eof = z.JumpToClass("img", "small-pic")
	s.Require().True(eof)
}

func (s *testSuite) TestJumpToTag() {
	z := NewTokenizer(strings.NewReader(rawHTML))

	attr, eof := z.JumpToTag("a")
	s.Require().False(eof)
	s.Require().Len(attr, 1)
	s.Require().Equal("#test", attr["href"])

	// Not found until EOF.
	_, eof = z.JumpToTag("a")
	s.Require().True(eof)
}

func (s *testSuite) TestExpandToken() {
	z := NewTokenizer(strings.NewReader(rawHTML))

	_, eof := z.JumpToTag("li")
	s.Require().False(eof)

	buf, eof := z.ExpandToken()
	s.Require().False(eof)

	dom := `
	<ol>option 1</ol>
	<ol>option 2</ol>
	<ol>option 3</ol>
`
	s.Require().Equal(dom, buf.String())
}

func TestFunc(t *testing.T) {
	suite.Run(t, new(testSuite))
}
