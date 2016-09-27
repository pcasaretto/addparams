package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func main() {
	err := Doit(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func Doit(r io.Reader, w io.Writer) error {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if err := z.Err(); err != io.EOF {
				return err
			}
			return nil
		case html.StartTagToken:
			tn, hasAttr := z.TagName()
			morAttr := true
			var key, val []byte
			if len(tn) == 1 && tn[0] == 'a' && hasAttr {
				fmt.Fprint(w, "<a")
				for morAttr {
					key, val, morAttr = z.TagAttr()
					if string(key) == "href" {
						newHREF := addParams(val)
						fmt.Fprintf(w, ` %s="%s"`, key, newHREF)
					} else {
						fmt.Fprintf(w, ` %s="%s"`, key, val)
					}
				}
				fmt.Fprint(w, ">")
				continue
			}
			fmt.Fprintf(w, "<%s", tn)
			if hasAttr {
				for morAttr {
					key, val, morAttr = z.TagAttr()
					fmt.Fprintf(w, " %s=%s ", key, val)
				}
			}
			fmt.Fprint(w, ">")
		default:
			w.Write(z.Raw())
		}
	}
}

func addParams(s []byte) string {
	input := string(s)
	url, err := url.Parse(input)
	if err != nil {
		return input
	}
	if scheme := url.Scheme; scheme == "http" || scheme == "https" {
		values := url.Query()
		values.Add("utm_blah", "blahblah")
		url.RawQuery = values.Encode()
		return url.String()
	}
	return input
}
