// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Translate uses the Google translate API from the command line to translate
// its arguments. By default it auto-detects the input language and translates
// to English.
// clone from "robpike/translate"

package main 

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"bytes"
)

var (
	// https://developers.google.com/translate/v2/using_rest#supported-query-params
	key    = flag.String("key", "", "Google API key (defaults to $GOOGLEAPIKEY)")
	target = flag.String("to", "en", "destination language (two-letter code)")
	source = flag.String("from", "", "source language (two-letter code); auto-detected by default")
)

type Response struct {
	Data struct {
		Translations []Translation
	}
}

type Translation struct {
	TranslatedText         string
	DetectedSourceLanguage string
}

func DoTrans(key string, target string, text string) string {
	var outmsg bytes.Buffer
	
	v := make(url.Values)
	v.Set("key", key)
	v.Set("target", target)
	v.Set("q", text)
	url := "https://www.googleapis.com/language/translate/v2?" + v.Encode()
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		log.Fatal(err)
	}
	for _, t := range r.Data.Translations {
		outmsg.WriteString(fmt.Sprintf("%s (%s)\n", html.UnescapeString(t.TranslatedText), t.DetectedSourceLanguage))
	}
	return outmsg.String()
}
