// StackTracer traces stacks using the sourcegraph api
package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func parseForm(req *http.Request, values ...string) (form url.Values, err error) {
	req.ParseForm()
	form = req.PostForm
	err = checkForm(form, values...)
	return
}

func checkForm(data url.Values, values ...string) error {
	for _, s := range values {
		if len(data[s]) == 0 {
			return errors.New(s + " not passed")
		}
	}
	return nil
}

func serveParse(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		v, err := parseForm(r, "trace")
		if err != nil {
			log.Println(err)
			return
		}
		output := parse(v["trace"][0])
		data, err := json.Marshal(output)
		if err != nil {
			log.Println(err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(data))
	}
}

var goTraceRE = regexp.MustCompile(`([^ ]*\.go):(\d+)`)
var baseURL = "https://sourcegraph.com/"

func parse(trace string) string {
	out := make([]string, 0)
	for _, line := range strings.Split(trace, "\n") {
		m := goTraceRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		// Only works with github right now.
		i := strings.Index(m[1], "github.com/")
		if i == -1 {
			continue
		}
		path := strings.SplitN(m[1][i:], "/", 4)
		if len(path) != 4 {
			continue
		}
		repo := strings.Join(path[0:3], "/")
		out = append(out, baseURL+repo+"/.tree/"+path[3]+
			"#startline="+m[2]+"&endline="+m[2])
	}
	if len(out) == 0 {
		return "no results"
	}
	return strings.Join(out, "\n")
}

func main() {
	http.HandleFunc("/parse", serveParse)
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("listening on localhost:8877")
	log.Fatal(http.ListenAndServe(":8877", nil))
}
