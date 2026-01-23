package enums

import (
	"database/sql/driver"
	"fmt"
)

type Indexing string

const (
	ALL               Indexing = "all"
	FOLLOW            Indexing = "follow"
	INDEX             Indexing = "index"
	INDEXIFEMBEDDED   Indexing = "indexifembedded"
	MAX_IMAGE_PREVIEW Indexing = "max-image-preview"
	MAX_SNIPPET       Indexing = "max-snippet"
	MAX_VIDEO_PREVIEW Indexing = "max-video-preview"
	NOAI              Indexing = "noai"
	NOARCHIVE         Indexing = "noarchive"
	NOCACHE           Indexing = "nocache"
	NOFOLLOW          Indexing = "nofollow"
	NOIMAGEAI         Indexing = "noimageai"
	NOIMAGEINDEX      Indexing = "noimageindex"
	NOINDEX           Indexing = "noindex"
	NOINDEXIFEMBEDDED Indexing = "noindexifembedded"
	NONE              Indexing = "none"
	NOODP             Indexing = "noodp" // Deprecated, kept for compatibility
	NOSNIPPET         Indexing = "nosnippet"
	NOTRANSLATE       Indexing = "notranslate"
	NOYDIR            Indexing = "noydir" // Deprecated, kept for compatibility
	UNAVAILABLE_AFTER Indexing = "unavailable_after"
)

func (i *Indexing) Scan(value interface{}) error {
	switch v := value.(type) {
	case nil:
		*i = ""
		return nil
	case string:
		*i = Indexing(v)
		return nil
	case []byte:
		*i = Indexing(string(v))
		return nil
	default:
		return fmt.Errorf("unsupported Scan type for Indexing: %T", value)
	}
}

func (i Indexing) Value() (driver.Value, error) {
	return string(i), nil
}

func (i Indexing) String() string {
	return string(i)
}
