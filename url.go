package url

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/purell"
)

type URL struct {
	Scheme    string
	Host      string
	Path      string
	Dir       string
	Base      string
	BaseQuery string
	RawQuery  string
	URLObject *url.URL
}

type ErrorURL struct {
	What string
}

func (e ErrorURL) Error() string {
	return fmt.Sprintf("URL: %s", e.What)
}

func StripWWW(host string) (string, bool) {
	if host[:4] == "www." {
		return host[4:], true
	}
	return host, false
}

func stripTrSlash(s string) string {
	if s[len(s)-1:] == "/" {
		s = s[:len(s)-1]
	}
	return s
}

func IsSubdomain(host string, u URL) bool {
	host, _ = StripWWW(host)
	host2, _ := StripWWW(u.Host)

	return strings.HasSuffix(host, host2)
}

func ParseFromObj(urlObj *url.URL) (*URL, error) {
	u := URL{Scheme: urlObj.Scheme,
		Host:      urlObj.Host,
		RawQuery:  urlObj.RawQuery,
		URLObject: urlObj}

	trSlash := false
	if urlObj.RawQuery != "" {
		if urlObj.RawQuery[len(urlObj.RawQuery)-1:] == "/" {
			u.RawQuery = stripTrSlash(urlObj.RawQuery)
			trSlash = true
		}
	}

	if urlObj.Path == "/" {
		u.Dir = "/"
		u.Path = "/"
		if u.RawQuery != "" {
			u.BaseQuery = "?" + u.RawQuery
			u.Path += u.BaseQuery
		}
		if trSlash {
			u.Path += "/"
		}
		return &u, nil
	}

	if urlObj.Path != "" {
		p := urlObj.Path
		if urlObj.Path[len(urlObj.Path)-1:] == "/" {
			p = urlObj.Path[:len(urlObj.Path)-1]
			trSlash = true
		}

		u.Dir = path.Dir(p)
		base := path.Base(p)

		if base != "" {
			u.Base = base
			u.BaseQuery = path.Base(p)
		}
		u.Path = u.Dir
		if u.Dir != "/" {
			u.Path += "/"
		}
	}

	if u.RawQuery != "" {
		u.BaseQuery += "?" + u.RawQuery
	}

	u.Path += u.BaseQuery

	if trSlash {
		u.Path += "/"
	}

	return &u, nil
}

func Parse(u string) (*URL, error) {
	nurl, err := purell.NormalizeURLString(u, purell.FlagsSafe)
	if err != nil {
		return nil, err
	}

	urlObj, err := url.Parse(nurl)

	if err != nil {
		return nil, err
	}

	return ParseFromObj(urlObj)
}

func (u *URL) String() string {
	if u.IsAbs() {
		return u.Scheme + "://" + u.Host + u.Path
	} else {
		return u.Path
	}
}

func (u *URL) ResolveReference(with *URL) (*URL, error) {
	to := u.URLObject.ResolveReference(with.URLObject)
	return ParseFromObj(to)
}

func (u *URL) IsAbs() bool {
	return u.URLObject.IsAbs()
}
