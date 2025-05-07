package crawler

// PageResult represents a single page with link count and SEO data
type PageResult struct {
	URL             string
	Count           int
	Title           string
	StatusCode      int
	WordCount       int
	InternalLinks   int
	ExternalLinks   int
	ImagesCount     int
	HasH1           bool
	H1Count         int
	HasMeta         bool
	HasCanonical    bool
	CanonicalURL    string
	MetaDescription string
}
