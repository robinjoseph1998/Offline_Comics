package utils

type Result struct {
	Month       string `json:"month"`
	Num         int    `json:"num"`
	Link        string `json:"link"`
	Year        string `json:"year"`
	News        string `json:"news"`
	SafeTitle   string `json:"safe_title"`
	Trasnscript string `json:"transcript"`
	Alt         string `json:"slt"`
	Img         string `json:"img"`
	Title       string `json:"title"`
	Day         string `json:"day"`
}

type Job struct {
	Number int
}
