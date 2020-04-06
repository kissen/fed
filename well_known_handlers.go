package main

import (
	"encoding/json"
	"gitlab.cs.fau.de/kissen/fed/config"
	"log"
	"net/http"
	"path"
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
