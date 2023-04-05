package selector

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

type Selector interface {
	SelectString(src string) string
	SelectStringList(src string) []string
}

type NodeSelector interface {
	Select(node *html.Node) string
	SelectList(node *html.Node) []string
	SelectNode(node *html.Node) *html.Node
	SelectNodeList(node *html.Node) []*html.Node
}

type LinkSelector struct{}

func (LinkSelector) Select(node *html.Node) string {
	panic("Select() can not apply to LinkSelector")
}

func (LinkSelector) SelectList(node *html.Node) []string {
	nodes := goquery.NewDocumentFromNode(node).Find("a").Nodes
	var links []string
	for _, node := range nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	return links
}

func (LinkSelector) SelectNode(node *html.Node) *html.Node {
	panic("SelectNode() can not apply to LinkSelector")
}

func (LinkSelector) SelectNodeList(node *html.Node) []*html.Node {
	panic("SelectNodeList() can not apply to LinkSelector")
}

type CssSelector struct {
	SelectorText string
	AttrName     string
	Pretty       bool
}

func (this CssSelector) Select(node *html.Node) string {
	nodes := goquery.NewDocumentFromNode(node).Find(this.SelectorText).Nodes
	if len(nodes) < 1 {
		return ""
	}
	return this.getValue(nodes[0])
}

func (this CssSelector) SelectList(node *html.Node) []string {
	nodes := goquery.NewDocumentFromNode(node).Find(this.SelectorText).Nodes
	var results []string
	if len(nodes) < 1 {
		return nil
	}
	for _, node := range nodes {
		results = append(results, this.getValue(node))
	}
	return results
}

func (this CssSelector) SelectNode(node *html.Node) *html.Node {
	if this.AttrName != "" {
		log.Warn("CssSelector.SelectNode() ignore AttrName")
	}
	nodes := goquery.NewDocumentFromNode(node).Find(this.SelectorText).Nodes
	if len(nodes) < 1 {
		return nil
	}
	return nodes[0]
}

func (this CssSelector) SelectNodeList(node *html.Node) []*html.Node {
	return goquery.NewDocumentFromNode(node).Find(this.SelectorText).Nodes
}

func (this CssSelector) getValue(node *html.Node) string {
	doc := goquery.NewDocumentFromNode(node)
	var result string
	var err error = nil
	switch {
	case strings.EqualFold("", this.AttrName) || strings.EqualFold("innerHtml", this.AttrName):
		result, err = doc.Html() //content exclude tag itself
	case strings.EqualFold("outerHtml", this.AttrName):
		result = getOuterHtml(doc) //content include tag itself
	case strings.EqualFold("text", this.AttrName):
		result = doc.Text()
	case strings.EqualFold("allText", this.AttrName):
		result = getAllText(doc)
	default:
		result, _ = doc.Attr(this.AttrName)
	}
	if err != nil {
		log.Errorf("failed to getValue, err：%+v", err)
		return ""
	}
	if this.Pretty {

	}
	return strings.TrimSpace(result)
}

func getOuterHtml(doc *goquery.Document) string {
	if len(doc.Parent().Nodes) < 1 {
		htmlStr, _ := doc.Html()
		return htmlStr
	}
	parentNode := doc.Parent().Get(0)
	htmlStr, _ := goquery.NewDocumentFromNode(parentNode).Html()
	return htmlStr
}

func getAllText(doc *goquery.Document) string {
	var allText string
	if len(doc.Children().Nodes) < 1 {
		return strings.Replace(goquery.NewDocumentFromNode(doc.Get(0)).Text(), "\n", "", -1)
	}
	for _, v := range doc.Children().Nodes {
		subStr := getAllText(goquery.NewDocumentFromNode(v))
		if subStr == "" {
			continue
		}
		allText += strings.TrimSpace(subStr) + " "
	}
	return allText
}

type RegexSelector struct {
	RegexStr string
	regexp   *regexp.Regexp
}

func NewRegexSelector(regexStr string) (*RegexSelector, error) {
	if regexStr == "" {
		return nil, errors.New("regexStr must not be empty")
	}
	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, err
	}
	regexSelector := RegexSelector{
		RegexStr: regexStr,
		regexp:   regex,
	}
	return &regexSelector, nil
}

func (this RegexSelector) SelectString(src string) string {
	return this.regexp.FindString(src)
}

func (this RegexSelector) SelectStringList(src string) []string {
	return this.regexp.FindAllString(src, -1)
}

type ReplaceSelector struct {
	RegexStr    string
	regexp      *regexp.Regexp
	Replacement string
}

func NewReplaceSelector(regexStr string, replacement string) (*ReplaceSelector, error) {
	if regexStr == "" {
		return nil, errors.New("expr must not be empty")
	}
	_regexp, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, err
	}
	replaceSelector := ReplaceSelector{
		RegexStr:    regexStr,
		regexp:      _regexp,
		Replacement: replacement,
	}
	return &replaceSelector, nil
}

func (this ReplaceSelector) SelectString(text string) string {
	return this.regexp.ReplaceAllString(text, this.Replacement)
}

func (this ReplaceSelector) SelectStringList(regexStr string) []string {
	panic("SelectStringList() can not apply to ReplaceSelector")
}
