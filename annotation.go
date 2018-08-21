package main

type Annotation struct {
	Context    string   `json:"@context"`
	Type       string   `json:"@type"`
	ID         string   `json:"@id"`
	Motivation []string `json:"motivation"`
	On         []struct {
		Type     string `json:"@type"`
		Full     string `json:"full"`
		Selector struct {
			Type    string `json:"@type"`
			Default struct {
				Type  string `json:"@type"`
				Value string `json:"value"`
			} `json:"default"`
			Item struct {
				Type  string `json:"@type"`
				Value string `json:"value"`
			} `json:"item"`
		} `json:"selector"`
		Within struct {
			ID   string `json:"@id"`
			Type string `json:"@type"`
		} `json:"within"`
	} `json:"on"`
	Resource []struct {
		Type  string `json:"@type"`
		Chars string `json:"chars"`

		Format string `json:"format"`
	} `json:"resource"`
}

func (a *Annotation) Manifest() string {
	return a.On[0].Within.ID
}

func (a *Annotation) Canvas() string {
	return a.On[0].Full
}

var AnnotationListTemplate = `
{
	"@context": "http://iiif.io/api/presentation/2/context.json",
	"@id": "http://annotation-studio.loc/annotation/list/1",
	"@type": "sc:AnnotationList",
  
	"resources": [
	    {{ StringsJoin . ", " }}
	]
  }
`
