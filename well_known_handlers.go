package main

import (
	"encoding/json"
	"fmt"
	"gitlab.cs.fau.de/kissen/fed/config"
	"gitlab.cs.fau.de/kissen/fed/fedcontext"
	"gitlab.cs.fau.de/kissen/fed/fediri"
	"log"
	"net/http"
	"path"
	"strings"
)

// GET /.well-known/nodeinfo
func GetNodeInfo(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNodeInfo()")

	p := *config.Get().Base
	p.Path = path.Join(p.Path, ".well-known", "nodeinfo", "2.0.json")

	href := p.String()
	rel := "http://nodeinfo.diaspora.software/ns/schema/2.0"

	reply := map[string]interface{}{
		"links": map[string]interface{}{
			"href": href,
			"rel":  rel,
		},
	}

	ReplyWithJSON(w, r, reply)
}

// GET /.well-known/nodeinfo/2.0.json
func GetNodeInfo20(w http.ResponseWriter, r *http.Request) {
	log.Println("GetNodeInfo()")

	reply := map[string]interface{}{
		"version": "2.0",
		"software": map[string]interface{}{
			"name":    "fed",
			"version": "0.x",
		},
		"services": map[string]interface{}{
			"inbound":  []string{},
			"outbound": []string{},
		},
		"protocols": []string{
			"activitypub",
		},
		"openRegistrations": false,
		"metadata": map[string]interface{}{
			"accountActivationRequired": false,
			"features":                  []string{},
			"federation": map[string]interface{}{
				"enabled":    true,
				"exclusions": false,
				"mrf_policies": []string{
					"NoOpPolicy",
				},
				"quarantined_instances": []string{},
			},
			"invitesEnabled": false,
			"mailerEnabled":  false,
			"nodeDescrption": "development instance",
			"nodeName":       config.Get().Base,
			"private":        false,
		},
	}

	ReplyWithJSON(w, r, reply)
}

// GET /.well-known/webfinger?resource=...
func GetWebfinger(w http.ResponseWriter, r *http.Request) {
	storage := fedcontext.Context(r).Storage
	configuration := config.Get()

	resource, ok := FormValue(r, "resource")
	if !ok {
		ApiError(w, r, "missing resource", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(resource, "acct:") {
		ApiError(w, r, "only acct supported", http.StatusNotImplemented)
		return
	}

	qs := strings.TrimPrefix(resource, "acct:")

	q := strings.Split(qs, "@")
	if len(q) != 2 {
		ApiError(w, r, "acct needs to have format username@server", http.StatusBadRequest)
		return
	}

	username := q[0]
	hostname := q[1]

	if hostname != configuration.Base.Hostname() {
		msg := fmt.Sprintf("bad hostname got=%v expected=%v", hostname, configuration.Base.Hostname())
		ApiError(w, r, msg, http.StatusBadRequest)
		return
	}

	if _, err := storage.RetrieveUser(username); err != nil {
		ApiError(w, r, err, http.StatusNotFound)
		return
	}

	href := fediri.ActorIRI(username).String()

	reply := map[string]interface{}{
		"subject": fmt.Sprintf("acct:%v@%v", username, hostname),
		"links": []interface{}{
			map[string]interface{}{
				"href": href,
				"rel":  "http://webfinger.net/rel/profile-page",
				"type": "text/html",
			},
			map[string]interface{}{
				"href": href,
				"rel":  "self",
				"type": "application/activity+json",
			},
			map[string]interface{}{
				"href": href,
				"rel":  "self",
				"type": AP_TYPE,
			},
		},
	}

	ReplyWithJSON(w, r, reply)
}

// GET /.well-known/host-meta
func GetHostMeta(w http.ResponseWriter, r *http.Request) {
	format := `<?xml version="1.0" encoding="UTF-8"?>` + "\n" +
		`<XRD xmlns="http://docs.oasis-open.org/ns/xri/xrd-1.0">` + "\n" +
		`  <Link rel="lrdd" type="application/xrd+xml" template="https://%s/.well-known/webfinger?resource={uri}"/>` + "\n" +
		`</XRD>` + "\n"

	xml := fmt.Sprintf(format, config.Get().Base.Host)
	bs := []byte(xml)

	w.Header().Add("Content-Type", "application/xrd+xml; charset=utf-8")

	if _, err := w.Write(bs); err != nil {
		log.Printf("writing xml to client failed: %v", err)
	}
}

func ReplyWithJSON(w http.ResponseWriter, r *http.Request, m map[string]interface{}) {
	bs, err := json.Marshal(m)
	if err != nil {
		ApiError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if _, err := w.Write(bs); err != nil {
		log.Printf("writing json to client failed: %v", err)
	}
}
