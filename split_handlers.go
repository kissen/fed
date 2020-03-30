package main

import (
	"net/http"
)

// This file contains various meta handlers that call
// either the web interface or the api (or an error).
//
// We should really do this with gorilla mux...

// GET /
func GetIndex(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebGetIndex(w, r)
	}
}

// GET|POST /user/outbox
func GetPostOutbox(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		ApGetPostOutbox(w, r)
	} else {
		GetRemote(w, r)
	}
}

// GET|POST /user/inbox
func GetPostInbox(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		ApGetPostInbox(w, r)
	} else {
		GetRemote(w, r)
	}
}

// GET /liked
func GetLiked(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebGetLiked(w, r)
	}
}

// GET /following
func GetFollowing(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebGetFollowing(w, r)
	}
}

// GET /followers
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebGetFollowers(w, r)
	}
}

// GET /login
func GetLogin(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebGetLogin(w, r)
	}
}

// POST /login
func PostLogin(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebPostLogin(w, r)
	}
}

// POST /logout
func PostLogout(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebPostLogout(w, r)
	}
}

// GET /storage/*
func ApGetStorage(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		ApGetPostActivity(w, r)
	} else {
		Error(w, r, http.StatusNotImplemented)
	}
}

// GET /remote/*
func GetRemote(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		ApGetRemote(w, r)
	} else {
		WebGetRemote(w, r)
	}
}

// POST /submit
func PostSubmit(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebPostSubmit(w, r)
	}
}

// POST /reply
func PostReply(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebPostReply(w, r)
	}
}

// POST /repeat
func PostRepeat(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebPostRepeat(w, r)
	}
}

// POST /like
func PostLike(w http.ResponseWriter, r *http.Request) {
	if IsActivityPubRequest(r) {
		WrongContentType(w, r)
	} else {
		WebPostLike(w, r)
	}
}
