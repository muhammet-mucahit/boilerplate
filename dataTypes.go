package boilerplate

// FormData ...
type FormData struct {
	URLs []string `form:"urls" json:"urls" xml:"urls" binding:"required"`
}

// ResultFormData ..
type ResultFormData struct {
	URL    string `json:"url"`
	Error  error  `json:"-"`
	Result Result `json:"result"`
}

// Result ..
type Result struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	H1          string `json:"h1"`
	Content     string `json:"content"`
	WordCount   int    `json:"word_count"`
}
