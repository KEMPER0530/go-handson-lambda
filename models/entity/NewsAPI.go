package entity

type NewsAPI struct {
	Status       string     `json:"status"`
	TotalResults int8       `json:"totalResults"`
	Articles     []struct { // <- 構造体の中にネストさせて構造体を定義
		Source struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"source"`
		Author      string `json:author`
		Title       string `json:title`
		Description string `json:description`
		Url         string `json:url`
		UrlToImage  string `json:urlToImage`
		PublishedAt string `json:publishedAt`
		Content     string `json:content`
	} `json:"articles"`
}
