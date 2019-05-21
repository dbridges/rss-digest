package main

import (
	"regexp"
	"time"
	"html"
	"html/template"
	"encoding/xml"
)

type Feed struct {
	ID      string  `xml:"id"`
	Title   string  `xml:"title"`
	Links   []Link  `xml:"link"`
	Entries []Entry `xml:"entry"`
	Updated atomTime `xml:"updated"`
}

type Entry struct {
	ID      string   `xml:"id"`
	Title   string   `xml:"title"`
	Links   []Link   `xml:"link"`
	Summary string   `xml:"summary"`
	Content Content  `xml:"content"`
	Updated atomTime `xml:"updated"`
	Author  Author	 `xml:"author"`
}

func (f *Feed) FilteredEntries() []Entry {
	cutoff := time.Now().Add(-24*time.Hour)
	entries := make([]Entry, 0, len(f.Entries))
	for _, e := range f.Entries {
		if time.Time(e.Updated).After(cutoff) {
			entries = append(entries, e)
		}
	}
	return entries
}

type Content struct {
	Type string `xml:"type,attr"`
	BaseURI string `xml:"base,attr"`
	Body string `xml:",innerxml"`
}

func (c Content) HTML() template.HTML {
	// Remove CDATA tags
	cdata := regexp.MustCompile(`(?s)<!\[CDATA\[(.*?)\]\]>`)
	body := cdata.ReplaceAllString(c.Body, "$1")
	// Make path to images absolute
	pat := regexp.MustCompile(`<img src="\/`)
	body = pat.ReplaceAllString(body, "<img src=\"" + c.BaseURI + "/")
	return template.HTML(html.UnescapeString(body))
}

type Link struct {
	HREF string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type Author struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type atomTime time.Time

func (a *atomTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return err
	}
	*a = atomTime(parsed)
	return nil
}

func (a *atomTime) Time() time.Time {
	return time.Time(*a)
}

func (a *atomTime) LocalString() string {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return a.Time().String()
	}
	return a.Time().In(loc).Format("Jan 02, 2006 3:04PM")
}
