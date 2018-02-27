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
	html := string(bhtml)

	mux := func(ctx *fasthttp.RequestCtx) {
		split := strings.Split(string(ctx.Path()), "/")

		var user string
		var ns string
		var body = html

		if len(split) > 2 {
			user = split[1]
			ns = split[2]

			defaultTiddlers, err := getDefaultTiddlers(user, ns)
			if err == nil {
				body = strings.Replace(body, `GettingStarted\n`, defaultTiddlers, 1)
			}

			body = addTiddler(body, "$:/plugins/fiatjaf/remoteStorage/readonly", "yes")
			body = addTiddler(body, "$:/plugins/fiatjaf/remoteStorage/userAddress", user)
		}

		ctx.SetContentType("text/html")
		ctx.SetBody([]byte(body))
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

			return strings.Join(BLANK.Split(tid.Text, -1), " "), nil
		}
	}

	return "", errors.New("No remoteStorage support on " + user)
}

func addTiddler(html, key, value string) string {
	tiddler := `
    <div author="tiddly.alhur.es" title="` + key + `" type="text/vnd.tiddlywiki">
      <pre>` + value + `</pre>
    </div>
	`

	loc := STOREAREAINIT.FindStringIndex(html)
	return html[0:loc[1]] + tiddler + html[loc[1]:]
}

var STOREAREAINIT = regexp.MustCompile(`id=['"]storeArea['"][^>]+>`)
var BLANK = regexp.MustCompile("[\n\r \t]+")
