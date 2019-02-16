package main

import (
	"gogoscrapy/src"
	"gogoscrapy/src/downloader"
	"gogoscrapy/src/entity"
	"gogoscrapy/src/pipeline"
	"strings"
	"sunteng/commons/log"
	"time"
)

type DoubanMovieProc struct {
}

//`Process` process the page and fetch the info we want.
func (this *DoubanMovieProc) Process(page entity.IPage) error {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("failed to process page, url:%s, err:%+v", page.GetUrl().Text(), err)
		}
	}()
	var reqs []entity.IRequest
	for _, node := range page.GetHtmlNode().Links().Regex(`https://movie.douban.com/subject/[0-9]+/`).Nodes() {
		url := strings.Split(node.Text(), "#")[0]
		if strings.HasPrefix(url, "https://movie.douban.com/subject/") {
			req := entity.NewRequest(url).SetUseProxy(false)
			reqs = append(reqs, req)
		}
	}
	url := strings.Split(page.GetUrl().Text(), "?")[0]
	page.StoreField("url", url)
	page.AddTargetRequests(reqs...)

	if !strings.HasPrefix(url, "https://movie.douban.com/subject") {
		return nil
	}

	//only /subject/ has the content we want
	page.StoreField("mid", url[len("https://movie.douban.com/subject/"):len(url)-1])
	for _, node := range page.GetHtmlNode().Css("#content > h1 > span:nth-child(1)").Nodes() {
		page.StoreField("title", node.Text())
	}
	for _, node := range page.GetHtmlNode().Css("#interest_sectl > div.rating_wrap.clearbox > div.rating_self.clearfix > strong").Nodes() {
		page.StoreField("rating", node.Text())
	}
	for _, node := range page.GetHtmlNode().Css("#interest_sectl > div.rating_wrap.clearbox > div.rating_self.clearfix > div > div.rating_sum > a > span").Nodes() {
		page.StoreField("ratingPeople", node.Text())
	}
	for _, node := range page.GetHtmlNode().Css("#info").Nodes() {
		info := node.Text()
		for _, line := range strings.Split(info, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "导演:") {
				page.StoreField("director", line[len("导演: "):])
			}
			if strings.HasPrefix(line, "编剧:") {
				page.StoreField("scriptwriter", line[len("编剧: "):])
			}
			if strings.HasPrefix(line, "主演:") {
				page.StoreField("actor", line[len("主演: "):])
			}
			if strings.HasPrefix(line, "类型:") {
				page.StoreField("type", line[len("类型: "):])
			}
			if strings.HasPrefix(line, "制片国家/地区:") {
				page.StoreField("country", line[len("制片国家/地区: "):])
			}
			if strings.HasPrefix(line, "语言:") {
				page.StoreField("language", line[len("语言: "):])
			}
			if strings.HasPrefix(line, "上映日期:") {
				page.StoreField("releaseDate", line[len("上映日期: "):])
			}
			if strings.HasPrefix(line, "IMDb链接:") {
				page.StoreField("imdb", line[len("IMDb链接: "):])
			}
		}
	}
	return nil
}

func main() {
	gogodjzhuSpider := src.NewSpider(&DoubanMovieProc{})
	gogodjzhuSpider.Downloader(downloader.NewSimpleDownloader(10*time.Second, nil))
	gogodjzhuSpider.Pipeline(pipeline.NewConsolePipeline())
	gogodjzhuSpider.DownloadCoroutineNum(3)
	gogodjzhuSpider.DownloadInterval(5 * time.Second)
	gogodjzhuSpider.AddStartUrl("https://movie.douban.com/", "https://movie.douban.com/chart")
	gogodjzhuSpider.Start()
}
