package ghubic

/*
   {
       "bytes": 529012,
       "content_type": "application/octet-stream",
       "hash": "4ffe4f4b940f27082b04fb4435c097f5",
       "last_modified": "2013-08-12T22:34:02.082030",
       "name": "~WRD3993.tmp"
   }
*/

type Files []File

type File struct {
	Bytes        int64  `json:"bytes"`
	ContentType  string `json:"content_type"`
	Hash         string `json:"hash"`
	LastModified string `json:"last_modified"`
	Name         string `json:"name"`
}
