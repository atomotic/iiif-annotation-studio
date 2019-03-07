package main

import "encoding/json"

// Annotation struct
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

// Manifest of the annotation
func (a *Annotation) Manifest() string {
	return a.On[0].Within.ID
}

// Canvas of the annotation
func (a *Annotation) Canvas() string {
	return a.On[0].Full
}

// AnnotationList is a list of annotation
type AnnotationList struct {
	Context   string            `json:"@context"`
	ID        string            `json:"@id"`
	Type      string            `json:"@type"`
	Resources []json.RawMessage `json:"resources"`
}
