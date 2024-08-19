package singlepagify

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"slices"
	"strings"
)

func Process(dir, file string) (*html.Node, error) {
	f, err := os.Open(dir + string(os.PathSeparator) + file)
	if err != nil {
		return nil, err
	}
	n, err := html.Parse(f)
	if err != nil {
		return nil, err
	}
	ns := elementsByTags(n, "link", "script")
	for i := range ns {
		processLink(ns[i], dir)
		processScript(ns[i], dir)
	}
	return n, nil
}

func elementsByTags(n *html.Node, tags ...string) []*html.Node {
	if n == nil {
		return nil
	}
	var out []*html.Node
	switch n.Type {
	case html.ErrorNode:
	case html.TextNode:
	case html.DocumentNode:
		if n.FirstChild != nil {
			for _, x := range elementsByTags(n.FirstChild, tags...) {
				out = append(out, x)
			}
		}
	case html.ElementNode:
		if slices.Contains(tags, n.Data) {
			out = append(out, n)
		}
		if n.FirstChild != nil {
			for _, x := range elementsByTags(n.FirstChild, tags...) {
				out = append(out, x)
			}
		}
	case html.CommentNode:
	case html.DoctypeNode:
		if n.FirstChild != nil {
			for _, x := range elementsByTags(n.FirstChild, tags...) {
				out = append(out, x)
			}
		}
	}

	for _, x := range elementsByTags(n.NextSibling, tags...) {
		out = append(out, x)
	}
	return out
}

func processLink(n *html.Node, dir string) {
	if n == nil || n.Data != "link" {
		return
	}
	var ok bool
	var fn string
	for i := 0; !(ok && fn != "") && i < len(n.Attr); i += 1 {
		if n.Attr[i].Key == "rel" {
			if ok = n.Attr[i].Val == "stylesheet"; !ok {
				return
			}
		}
		if n.Attr[i].Key == "href" {
			fn = n.Attr[i].Val
		}
	}
	if !ok || fn == "" {
		return
	}
	n.Data = "style"
	n.Attr = nil
	fn = dir + string(os.PathSeparator) + fn
	bs, err := os.ReadFile(fn)
	if err != nil {
		fmt.Printf("could not read %q, skipping\n", fn)
		return
	}
	t := &html.Node{
		Type: html.TextNode,
		Data: strings.TrimSpace(string(bs)),
	}
	n.AppendChild(t)
}

func processScript(n *html.Node, dir string) {
	if n == nil || n.Data != "script" {
		return
	}
	var fn string
	for i := 0; fn == "" && i < len(n.Attr); i += 1 {
		if n.Attr[i].Key == "src" {
			fn = n.Attr[i].Val
		}
	}
	if fn == "" {
		return
	}
	fn = dir + string(os.PathSeparator) + fn
	bs, err := os.ReadFile(fn)
	if err != nil {
		fmt.Printf("could not read %q, skipping\n", fn)
		return
	}
	n.Attr = []html.Attribute{{Key: "type", Val: "application/javascript"}}
	t := &html.Node{
		Type: html.TextNode,
		Data: strings.TrimSpace(string(bs)),
	}
	n.AppendChild(t)
}
