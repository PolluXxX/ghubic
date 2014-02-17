package ghubic

type Files []File

type File struct {
	Bytes        int64  `json:"bytes"`
	ContentType  string `json:"content_type"`
	Hash         string `json:"hash"`
	LastModified string `json:"last_modified"`
	Name         string `json:"name"`
}
