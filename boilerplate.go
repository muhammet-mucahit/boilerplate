package boilerplate

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jlubawy/go-boilerpipe"
)

// ContentFinder ...
type ContentFinder struct{}

func getH1(document *goquery.Document) string {
	h1 := ""

	titleElement := document.Find("h1")
	if titleElement != nil && titleElement.Size() > 0 {
		h1 = titleElement.Text()
	}

	return strings.TrimSpace(h1)
}

func getWordsCount(doc *boilerpipe.Document) int {
	numWordsOfContent := 0
	for _, v := range doc.TextBlocks {
		if v.IsContent {
			numWordsOfContent += v.NumWords
		}
	}
	return numWordsOfContent
}

// getMetaContentWithSelector returns the content attribute of meta tag matching the selector
func getMetaContentWithSelector(document *goquery.Document, selector string) string {
	selection := document.Find(selector)
	content, _ := selection.Attr("content")
	return strings.TrimSpace(content)
}

func asd(res *http.Response) (string, string) {
	document, _ := goquery.NewDocumentFromReader(res.Body)
	return getH1(document), getMetaContentWithSelector(document, "meta[name#=(?i)^description$]")
}

func boilerplate(url string, ch chan *ResultFormData) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	userAgent := "Mozilla/5.0 (Windows NT 6.2; WOW64; rv:21.0) Gecko/20130514 Firefox/21.0"
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	h1, desc := asd(res) //"", ""

	doc, err := boilerpipe.ParseDocument(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	boilerpipe.ArticlePipeline.Process(doc)
	title := doc.Title
	wordCount := getWordsCount(doc)

	ch <- &ResultFormData{
		URL:   url,
		Error: err,
		Result: Result{
			Title:       title,
			Description: desc,
			H1:          h1,
			Content:     doc.Content(),
			WordCount:   wordCount,
		},
	}

	if err == nil && res != nil && res.StatusCode == http.StatusOK {
		res.Body.Close()
	}
}

// Find ...
func (cf *ContentFinder) Find(urls []string) []*ResultFormData {
	ch := make(chan *ResultFormData)
	responses := []*ResultFormData{}

	for _, url := range urls {
		go boilerplate(url, ch)
	}

	for {
		select {
		case r := <-ch:
			fmt.Printf("%s was fetched\n", r.URL)
			// fmt.Println(r.result)
			responses = append(responses, r)
			if len(responses) == len(urls) {
				return responses
			}
		default:
			// fmt.Printf(".")
			// time.Sleep(50 * time.Millisecond)
		}
	}
}
