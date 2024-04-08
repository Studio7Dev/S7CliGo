package GPT_CLI

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Href  string `json:"href"`
	Desc  string `json:"desc"`
	Title string `json:"title"`
}

type SearchEngine struct {
	client *http.Client
	query  string
}

func NewSearchEngine() *SearchEngine {
	return &SearchEngine{
		client: &http.Client{},
		query:  "golang download",
	}
}

func (engine *SearchEngine) encodeURLParams(params map[string]string) string {
	data := url.Values{}
	for key, value := range params {
		data.Add(key, value)
	}
	return data.Encode()
}

func (engine *SearchEngine) fetchArticles() ([]*Article, error) {
	const baseUrl = "https://priv.au/search"
	params := map[string]string{
		"q":                engine.query,
		"category_general": "1",
		"language":         "en",
		"time_range":       "",
		"safesearch":       "0",
		"theme":            "simple",
	}
	data := engine.encodeURLParams(params)
	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed creating request: %w", err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", `categories=general; language=en; locale=en; autocomplete=google; image_proxy=1; method=POST; safesearch=1; theme=simple; results_on_new_tab=0; doi_resolver=oadoi.org; simple_style=macchiato; center_alignment=0; advanced_search=0; query_in_title=1; infinite_scroll=1; search_on_category_select=1; hotkeys=default; disabled_engines=; enabled_engines="artic__images\054ask__general\054bing__general\054bing images__images\054bpb__general\054openverse__images\054crowdview__general\054yep__general\054yep images__images\054curlie__general\054currency__general\054bahnhof__general\054deviantart__images\054ddg definitions__general\054wikidata__general\054duckduckgo__general\054duckduckgo images__images\054tineye__general\0541x__images\054flickr__images\054frinkiac__images\054material icons__images\054imgur__images\054library of congress__images\054lingva__general\054mozhi__general\054mwmbl__general\054pinterest__images\054presearch__general\054presearch images__images\054presearch videos__general\054qwant__general\054qwant images__images\054startpage__general\054tagesschau__general\054unsplash__images\054yahoo__general\054wiby__general\054alexandria__general\054wikibooks__general\054wikiquote__general\054wikisource__general\054wikispecies__general\054wikiversity__general\054wikivoyage__general\054wikicommons.images__images\054wolframalpha__general\054dictzone__general\054seznam__general\054mojeek__general\054naver__general\054yacy__general\054yacy images__images\054seekr images__images\054stract__general\054svgrepo__images\054wallhaven__images\054wikimini__general\054brave__general\054brave.images__images\054goo__general"; disabled_plugins=; enabled_plugins=; tokens=; maintab=on; enginetab=on; preferences=`)
	req.Header.Set("Origin", "null")
	req.Header.Set("Sec-CH-UA", `"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := engine.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed sending request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyText)))
	if err != nil {
		return nil, fmt.Errorf("failed parsing HTML: %w", err)
	}

	var articles []*Article
	doc.Find("article").Each(func(_ int, articleSel *goquery.Selection) {
		link := articleSel.Find("a")
		description := articleSel.Find("p")
		title := articleSel.Find("h3")

		article := &Article{
			Href:  link.AttrOr("href", ""),
			Desc:  description.Text(),
			Title: title.Text(),
		}
		articles = append(articles, article)
	})

	return articles, nil
}

func (engine *SearchEngine) Run(query string) string {
	engine.query = query
	articles, err := engine.fetchArticles()
	if err != nil {
		log.Fatalf("Error while running search: %v\n", err)
	}

	jsonData, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		fmt.Println(err)
		return "error"
	}

	return string(jsonData)
}
