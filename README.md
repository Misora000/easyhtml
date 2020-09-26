# Easy HTML Crawler

A easy-to-use HTML crawler based on `golang.org/x/net/html`.  
Find what you want intuitively, such as:
```
JumpToTag(tag string) (a Attrs, eof bool)
JumpToID(tag, id string) (a Attrs, eof bool)
JumpToClass(tag, class string) (a Attrs, eof bool)
ExpandToken() (*bytes.Buffer, bool)
```