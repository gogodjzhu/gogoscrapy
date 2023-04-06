package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/gogodjzhu/gogoscrapy"
	"github.com/gogodjzhu/gogoscrapy/downloader"
	entity2 "github.com/gogodjzhu/gogoscrapy/entity"
	utils2 "github.com/gogodjzhu/gogoscrapy/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DoubanBookProc struct {
}

//`Process` process the page and fetch the info we want.
func (this *DoubanBookProc) Process(page entity2.IPage) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("failed to process page, url:%s\nstackTrace:%s\nerr:%+v", page.GetUrl().Text(), utils2.GetStackTrace(), err)
		}
	}()
	var reqs []entity2.IRequest
	for _, node := range page.GetHtmlNode().Links().Regex(`https://book.douban.com/subject/[0-9]+/`).Nodes() {
		url := strings.Split(node.Text(), "#")[0]
		if strings.HasPrefix(url, "https://book.douban.com/subject/") {
			req, err := entity2.NewGetRequest(url)
			if err != nil {
				log.Errorf("failed to create request, url:%s, err:%+v", url, err)
				continue
			}
			req.SetUseProxy(true)
			//req.SetPriority(3)//set the priority if you want, greater will be processed first
			reqs = append(reqs, req)
		}
	}
	url := strings.Split(page.GetUrl().Text(), "?")[0]
	page.StoreField("url", url)
	page.AddTargetRequests(reqs...)

	if !strings.HasPrefix(url, "https://book.douban.com/subject") {
		return nil
	}

	//only /subject/ has the content we want
	page.StoreField("mid", url[len("https://book.douban.com/subject/"):len(url)-1])
	for _, node := range page.GetHtmlNode().Css("#wrapper > h1 > span").Nodes() {
		page.StoreField("title", node.Text())
	}
	if len(page.GetHtmlNode().Css("#wrapper > h1 > span").Nodes()) == 0 {
		return errors.New("page content invalid, html:" + page.GetRawText())
	}
	for _, node := range page.GetHtmlNode().Css("#interest_sectl > div > div.rating_self.clearfix > strong").Nodes() {
		page.StoreField("rating", node.Text())
	}
	for _, node := range page.GetHtmlNode().Css("#interest_sectl > div > div.rating_self.clearfix > div > div.rating_sum > span > a > span").Nodes() {
		page.StoreField("ratingPeople", node.Text())
	}
	for _, node := range page.GetHtmlNode().Css("#interest_sectl > div > div.rating_self.clearfix > div > div.rating_sum > span > a > span").Nodes() {
		page.StoreField("ratingPeople", node.Text())
	}
	for _, node := range page.GetHtmlNode().Css("#info").Nodes() {
		info := node.Html()
		for _, line := range strings.Split(info, "<br/>") {
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(line)))
			if err != nil {
				return err
			}
			line = doc.Text()
			line = strings.Replace(line, " ", "", -1)
			line = strings.Replace(line, "\n", "", -1)
			if strings.HasPrefix(line, "作者:") {
				page.StoreField("author", line[len("作者:"):])
			}

			if strings.HasPrefix(line, "出版社:") {
				page.StoreField("press", line[len("出版社:"):])
			}
			if strings.HasPrefix(line, "副标题:") {
				page.StoreField("subTitle", line[len("副标题:"):])
			}
			if strings.HasPrefix(line, "原作名:") {
				page.StoreField("originalTitle", line[len("原作名:"):])
			}
			if strings.HasPrefix(line, "译者:") {
				page.StoreField("translator", line[len("译者:"):])
			}
			if strings.HasPrefix(line, "出版年:") {
				page.StoreField("pressDate", line[len("出版年:"):])
			}
			if strings.HasPrefix(line, "页数:") {
				page.StoreField("pageNum", line[len("页数:"):])
			}
			if strings.HasPrefix(line, "定价:") {
				page.StoreField("price", line[len("定价:"):])
			}
			if strings.HasPrefix(line, "装帧:") {
				page.StoreField("binding", line[len("装帧:"):])
			}
			if strings.HasPrefix(line, "ISBN:") {
				page.StoreField("isbn", line[len("ISBN:"):])
			}
			if strings.HasPrefix(line, "从书:") {
				page.StoreField("series", line[len("从书:"):])
			}
		}
	}
	page.StoreField("tbl", "t_books")
	return nil
}

func main() {
	spider := gogoscrapy.NewSpider(&DoubanBookProc{})
	spider.Downloader(downloader.NewSimpleDownloader(10*time.Second, nil))
	spider.DownloadCoroutineNum(1)
	spider.DownloadInterval(5 * time.Second)
	spider.RetryTime(10)
	spider.AddStartUrl("https://book.douban.com/subject/27081847/")
	spider.Start()
}
