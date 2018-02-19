package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/ant0ine/go-webfinger"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

func main() {
	htmlfile, err := os.Open("index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer htmlfile.Close()
	bhtml, err := ioutil.ReadAll(htmlfile)
	if err != nil {
		log.Fatal(err)
	}

	mux := func(ctx *fasthttp.RequestCtx) {
		split := strings.Split(string(ctx.Path()), "/")

		var page string
		var user string
		var ns string
		switch len(split) {
		case 3:
			user = split[1]
			ns = split[2]
		case 4:
			user = split[1]
			ns = split[2]
			page = split[3]
		default:
			ctx.Redirect("https://tiddly.alhur.es/", 302)
			return
		}

		html := string(bhtml)

		defaultTiddlers, err := getDefaultTiddlers(user, ns)
		if err != nil {
			ctx.Error(err.Error(), 400)
			return
		}

		if page != "" {
			defaultTiddlers = page + " " + defaultTiddlers
		}

		html = setDefaultTiddlers(html, defaultTiddlers)

		ctx.SetContentType("text/html")
		ctx.SetBody([]byte(html))

	}

	log.Print("listening...")
	err = fasthttp.ListenAndServe(":"+os.Getenv("PORT"), mux)
	if err != nil {
		log.Print(err)
	}
}

func getDefaultTiddlers(user string, ns string) (string, error) {
	jrd, err := webfinger.Lookup(user, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to fetch webfinger JRD.")
	}

	for _, link := range jrd.Links {
		if link.Rel == "http://tools.ietf.org/id/draft-dejong-remotestorage" ||
			link.Rel == "remotestorage" ||
			link.Rel == "remoteStorage" {
			resp, err := http.Get(link.Href + "/public/tiddlers/" + ns + "/%2524%253A%252FDefaultTiddlers")
			if err != nil {
				return "", errors.Wrap(err, "failed to fetch $:/DefaultTiddlers from "+link.Href+"/public/tiddlers/"+ns+".")
			}

			defer resp.Body.Close()
			tid := struct {
				Text string `json:"text"`
			}{}
			err = json.NewDecoder(resp.Body).Decode(&tid)
			if err != nil {
				return "", errors.Wrap(err, "failed to parse $:/DefaultTiddlers JSON.")
			}

			return tid.Text, nil
		}
	}

	return "", errors.New("No remoteStorage support on " + user)
}

func setDefaultTiddlers(html string, names string) string {
	tiddler := `
    <div author="tiddly.alhur.es" title="$:/DefaultTiddlers" type="text/vnd.tiddlywiki">
      <pre>` + names + `</pre>
    </div>
	`

	loc := STOREAREAINIT.FindStringIndex(html)
	return html[0:loc[1]] + tiddler + html[loc[1]:]
}

var STOREAREAINIT = regexp.MustCompile(`id=['"]storeArea['"][^>]+>`)
