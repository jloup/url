package url

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

var testdata = []struct {
	in           string
	outURL       string
	outPath      string
	outDir       string
	outBase      string
	outBaseQuery string
}{
	{"http://www.host.com", "http://www.host.com", "", "", "", ""},
	{"http://www.host.com/", "http://www.host.com/", "/", "/", "", ""},
	{"http://host.com/", "http://host.com/", "/", "/", "", ""},
	{"http://www.host.com:80", "http://www.host.com", "", "", "", ""},
	{"http://www.host.com:80/", "http://www.host.com/", "/", "/", "", ""},
	{"http://host.com/root", "http://host.com/root", "/root", "/", "root", "root"},
	{"http://host.com:80/root", "http://host.com/root", "/root", "/", "root", "root"},
	{"http://host.com/root/document", "http://host.com/root/document", "/root/document", "/root", "document", "document"},
	{"http://host.com/root/subdir/document", "http://host.com/root/subdir/document", "/root/subdir/document", "/root/subdir", "document", "document"},
	{"http://host.com/root/subdir/document/", "http://host.com/root/subdir/document/", "/root/subdir/document/", "/root/subdir", "document", "document"},
	{"http://host.com/root/subdir/document?q=query", "http://host.com/root/subdir/document?q=query", "/root/subdir/document?q=query", "/root/subdir", "document", "document?q=query"},
	{"http://host.com?q=query", "http://host.com?q=query", "?q=query", "", "", "?q=query"},
	{"http://host.com?q=query/", "http://host.com?q=query/", "?q=query/", "", "", "?q=query"},
	{"http://host.com/?q=query", "http://host.com/?q=query", "/?q=query", "/", "", "?q=query"},
	{"http://host.com/?q=query/", "http://host.com/?q=query/", "/?q=query/", "/", "", "?q=query"},
	{"http://host.com/root/subdir/document?q=query/", "http://host.com/root/subdir/document?q=query/", "/root/subdir/document?q=query/", "/root/subdir", "document", "document?q=query"},
}

func TestBasic(t *testing.T) {

	for _, set := range testdata {

		p, err := Parse(set.in)

		if err != nil {
			t.Fatal(err)
		}

		if p.String() != set.outURL {
			t.Fatalf("ERROR FOR SET => in '%v' want '%v' got '%v'\n", set.in, set.outURL, p.String())
		}

		if p.String() != p.URLObject.String() {
			t.Fatalf("POSSIBLE BAD FORMATTING => for '%v' URL gives '%v' net.url gives '%v' (want '%v')\n", set.in, p.String(), p.URLObject.String(), set.outURL)
		}

		if p.Path != set.outPath {
			t.Fatalf("BAD PATH => in '%v' want '%v' got '%v'\n", set.in, p.Path, set.outPath)
		}

		if p.Dir != set.outDir {
			t.Fatalf("BAD DIR => in '%v' want '%v' got '%v'\n", set.in, p.Dir, set.outDir)
		}

		if p.Base != set.outBase {
			t.Fatalf("BAD BASE => in '%v' want '%v' got '%v'\n", set.in, p.Base, set.outBase)
		}

		if p.BaseQuery != set.outBaseQuery {
			t.Fatalf("BAD BASEQUERY => in '%v' want '%v' got '%v'\n", set.in, set.outBaseQuery, p.BaseQuery)
		}

		t.Logf("in '%v' => '%v'\n", set.in, p.String())
	}
}

func TestRelative(t *testing.T) {
	where, _ := Parse("http://www.dealerofpeopleemotions.com/yo/wordpress/actu")
	//where, _ := url.Parse("http://www.dealerofpeopleemotions.com/yo/wordpress/actu")

	url1, _ := Parse("/from/root/")
	url2, _ := Parse("from/current")
	url3, _ := Parse("../from/parent")

	//url1, _ := url.Parse("/from/root/")
	//url2, _ := url.Parse("./from/current")
	//url3, _ := url.Parse("../from/parent")

	fmt.Println(url1)
	u1R, _ := where.ResolveReference(url1)
	u2R, _ := where.ResolveReference(url2)
	u3R, _ := where.ResolveReference(url3)

	//u1R := where.ResolveReference(url1)
	//u2R := where.ResolveReference(url2)
	//u3R := where.ResolveReference(url3)

	t.Log(url1, "=>", u1R)
	t.Log(url2, "=>", u2R)
	t.Log(url3, "=>", u3R)
}

func isSubdomain(parent, sub *url.URL) bool {
	return strings.HasSuffix(sub.Host, parent.Host)
}

func TestSubDomain(t *testing.T) {
	url1, _ := url.Parse("http://lemonde.fr/")
	url2, _ := url.Parse("http://bigbrowser.blog.lemonde.fr")
	url3, _ := url.Parse("http://www.blog.lemonde.fr")

	t.Log(url2, "=>", isSubdomain(url1, url2))
	t.Log(url3, "=>", isSubdomain(url1, url3))

}
