package selector

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Selectable interface {
	Links() Selectable
	Regex(regex string) Selectable
	Html() string
	Text() string
	Match() bool
	Css(selectorText string) Selectable
	CssWithAttr(selectorText, attrName string) Selectable
	Replace(regex, replacement string) Selectable
	Nodes() []Selectable
}

type HtmlNode struct {
	Elements []*html.Node
}

func (hn *HtmlNode) Links() Selectable {
	var sourceTexts []string
	for _, elem := range hn.Elements {
		sourceTexts = append(sourceTexts, LinkSelector{}.SelectList(elem)...)
	}
	return PlainText{SourceTexts: sourceTexts}
}

func (hn *HtmlNode) Regex(regex string) Selectable {
	var sourceTexts []string
	regexSelector, err := NewRegexSelector(regex)
	if err != nil {
		panic(err)
	}
	for _, elem := range hn.Elements {
		elemStr, err := goquery.NewDocumentFromNode(elem).Html()
		if err != nil {
			panic(err)
		}
		sourceTexts = append(sourceTexts, regexSelector.SelectStringList(elemStr)...)
	}
	return PlainText{SourceTexts: sourceTexts}
}

func (hn *HtmlNode) Html() string {
	if hn.Elements == nil {
		return ""
	}
	if htmlStr, err := goquery.NewDocumentFromNode(hn.Elements[0]).Html(); err != nil {
		panic(err)
	} else {
		return htmlStr
	}
}

func (hn *HtmlNode) Text() string {
	if hn.Elements == nil {
		return ""
	}
	var text string
	for _, v := range hn.Elements {
		text += goquery.NewDocumentFromNode(v).Text()
	}
	return text
}

func (hn *HtmlNode) Match() bool {
	return len(hn.Elements) > 0
}

func (hn *HtmlNode) Css(selectorText string) Selectable {
	var nodes []*html.Node
	selector := CssSelector{
		SelectorText: selectorText,
		AttrName:     "outerHtml",
	}
	for _, elem := range hn.Elements {
		nodes = append(nodes, selector.SelectNodeList(elem)...)
	}
	return &HtmlNode{Elements: nodes}
}

// has attribute, consider as plaintext
func (hn *HtmlNode) CssWithAttr(selectorText string, attrName string) Selectable {
	var sourceTexts []string
	selector := CssSelector{SelectorText: selectorText, AttrName: attrName}
	for _, elem := range hn.Elements {
		sourceTexts = append(sourceTexts, selector.SelectList(elem)...)
	}
	return PlainText{SourceTexts: sourceTexts}
}

func (hn *HtmlNode) Replace(src, replacement string) Selectable {
	var retStrings []string
	selector, err := NewReplaceSelector(src, replacement)
	if err != nil {
		panic(err)
	}
	for _, elem := range hn.Elements {
		htmlStr := goquery.NewDocumentFromNode(elem).Text()
		retStrings = append(retStrings, selector.SelectString(htmlStr))
	}
	return PlainText{SourceTexts: retStrings}
}

func (hn *HtmlNode) Nodes() []Selectable {
	var selectables []Selectable
	for _, element := range hn.Elements {
		selectables = append(selectables, &HtmlNode{Elements: []*html.Node{element}})
	}
	return selectables
}

type PlainText struct {
	SourceTexts []string
}

func (pt PlainText) Links() Selectable {
	panic("Links() can not apply to PlainText")
}

func (pt PlainText) Regex(regex string) Selectable {
	regexSelector, err := NewRegexSelector(regex)
	if err != nil {
		panic(err)
	}
	var resultArr []string
	for _, v := range pt.SourceTexts {
		resultArr = append(resultArr, regexSelector.SelectStringList(v)...)
	}
	return PlainText{SourceTexts: resultArr}
}

func (pt PlainText) Html() string {
	panic("Html() can not apply to PlainText")
}

func (pt PlainText) Text() string {
	if len(pt.SourceTexts) < 1 {
		return ""
	}
	var text string
	for _, v := range pt.SourceTexts {
		text = text + v
	}
	return text
}

func (pt PlainText) Match() bool {
	return len(pt.SourceTexts) > 0
}

func (pt PlainText) Css(selector string) Selectable {
	panic("Css() can not apply to PlainText")
}

func (pt PlainText) CssWithAttr(selector string, attrName string) Selectable {
	panic("CssWithAttr() can not apply to PlainText")
}

func (pt PlainText) Replace(regexStr, replacement string) Selectable {
	if selector, err := NewReplaceSelector(regexStr, replacement); err != nil {
		panic(err)
	} else {
		var sourceTexts []string
		for _, v := range pt.SourceTexts {
			sourceTexts = append(sourceTexts, selector.SelectString(v))
		}
		return PlainText{SourceTexts: sourceTexts}
	}
}

func (pt PlainText) Nodes() []Selectable {
	var selectables []Selectable
	for _, sourceText := range pt.SourceTexts {
		selectables = append(selectables, PlainText{SourceTexts: []string{sourceText}})
	}
	return selectables
}
