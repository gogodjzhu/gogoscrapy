package selector

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func getHtmlNode() HtmlNode {
	reader := strings.NewReader(htmlStr)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}
	htmlNode := HtmlNode{[]*html.Node{doc.Get(0)}}
	return htmlNode
}
func TestHtmlNode_Nodes(t *testing.T) {
	TestHtmlNode_Links(t)
}

func TestHtmlNode_Links(t *testing.T) {
	expected := 67
	actualArr := getHtmlNode().Links().Nodes()
	actual := len(actualArr)
	if expected != actual {
		t.Errorf("failed test @ TestHtmlNode_Links, expecteds:%d, actuals:%d", expected, actual)
	}
}

func TestHtmlNode_Html(t *testing.T) {
	getHtmlNode().Html()
}

func TestHtmlNode_Text(t *testing.T) {
	expected := "http://gogodjzhu.com/"
	actual := getHtmlNode().Links().Nodes()[0].Text()
	if expected != actual {
		t.Errorf("failed test @ TestHtmlNode_Text, expecteds:%s, actuals:%s", expected, actual)
	}
}

func TestHtmlNode_Match(t *testing.T) {
	expected := true
	actual := getHtmlNode().Css("div").Match()
	if expected != actual {
		t.Errorf("failed test @ TestHtmlNode_Match, expecteds:%t, actuals:%t", expected, actual)
	}
}

func TestHtmlNode_Replace(t *testing.T) {
	expected := "class读取hdfs集群的全部配置"
	actualArr := getHtmlNode().Css(".post-title").Replace("命令行", "class").Nodes()
	actual := strings.TrimSpace(actualArr[3].Text())
	if expected != actual {
		t.Errorf("failed test @ TestHtmlNode_Replace, expecteds:%s, actuals:%s", expected, actual)
	}
}

func TestHtmlNode_Css(t *testing.T) {
	expectedStr :=
		`《送你一颗子弹》——读后随笔
西班牙内战——当世界还年轻的时候
Practices in memorising the scales and the chords
命令行读取hdfs集群的全部配置
使用Instrumentation计算java对象大小
Astah Professional 7.2.0/1ff236 破解工具
Ours Samplus – Deep Inside
All of me – [The jazz real book series]
# 嘘！老猪教你如何冲浪~
很酷的一些东西`
	expectedArr := strings.Split(expectedStr, "\n")
	selectable := getHtmlNode().Css(".post-title > a").Nodes()
	for i, s := range selectable {
		if expectedArr[i] != s.Text() {
			t.Errorf("failed test @ TestHtmlNode_Css, expecteds:%s, actuals:%s", expectedArr[i], s.Text())
		}
	}
}

func TestHtmlNode_CssWithAttr(t *testing.T) {
	expectedStr :=
		`http://gogodjzhu.com/index.php/read/371/
http://gogodjzhu.com/index.php/something/365/
http://gogodjzhu.com/index.php/something/353/
http://gogodjzhu.com/index.php/code/tools/347/
http://gogodjzhu.com/index.php/code/basic/329/
http://gogodjzhu.com/index.php/code/tools/326/
http://gogodjzhu.com/index.php/music/resources/322/
http://gogodjzhu.com/index.php/music/resources/316/
http://gogodjzhu.com/index.php/something/305/
http://gogodjzhu.com/index.php/code/303/`
	expectedArr := strings.Split(expectedStr, "\n")
	actualArr := getHtmlNode().CssWithAttr(".post-title > a", "href").Nodes()
	for i, s := range actualArr {
		if expectedArr[i] != s.Text() {
			t.Errorf("failed test @ TestHtmlNode_CssWithAttr, expecteds:%s, actuals:%s", expectedArr[i], s.Text())
		}
	}
}

func TestHtmlNode_Regex(t *testing.T) {
	expected := 3
	actualArr := getHtmlNode().CssWithAttr(".post-title > a", "href").Regex(".*something.*").Nodes()
	actual := len(actualArr)
	if expected != actual {
		t.Errorf("failed test @ TestHtmlNode_Regex, expecteds:%d, actuals:%d", expected, actual)
	}
}

func TestPlainText_Html(t *testing.T) {
	expected := "Html() can not apply to PlainText"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestPlainText_Html, should not be implemented")
		}
	}()
	getHtmlNode().Regex("http").Html()
}

func TestPlainText_Text(t *testing.T) {
	defer func() {
		if recover() != nil {
			t.Error("failed test @ TestPlainText_Text, panic")
		}
	}()
	getHtmlNode().Regex("http").Text()
}

func TestPlainText_Links(t *testing.T) {
	expected := "Links() can not apply to PlainText"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestPlainText_Links, should not be implemented")
		}
	}()
	getHtmlNode().Regex("http").Links()
}

func TestPlainText_Match(t *testing.T) {
	expected := true
	actual := getHtmlNode().Links().Match()
	if expected != actual {
		t.Errorf("failed test @ TestPlainText_Match, expecteds:%t, actuals:%t", expected, actual)
	}
}

func TestPlainText_Nodes(t *testing.T) {
	expected := 8
	actual := len(getHtmlNode().Regex("GoGo DJZhu").Nodes())
	if expected != actual {
		t.Errorf("failed test @ TestPlainText_Nodes, expecteds:%d, actuals:%d", expected, actual)
	}
}

func TestPlainText_Regex(t *testing.T) {
	TestPlainText_Nodes(t)
}

func TestPlainText_Replace(t *testing.T) {
	expected := "djzhu://gogodjzhu.com/"
	actual := getHtmlNode().Links().Nodes()[0].Replace("http", "djzhu").Text()
	if expected != actual {
		t.Errorf("failed test @ TestPlainText_Replace, expecteds:%s, actuals:%s", expected, actual)
	}
}

func TestPlainText_Css(t *testing.T) {
	expected := "Css() can not apply to PlainText"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestPlainText_Css, should not be implemented")
		}
	}()
	getHtmlNode().Links().Css("div")
}

func TestPlainText_CssWithAttr(t *testing.T) {
	expected := "CssWithAttr() can not apply to PlainText"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestPlainText_Css, should not be implemented")
		}
	}()
	getHtmlNode().Links().CssWithAttr("div", "class")
}
