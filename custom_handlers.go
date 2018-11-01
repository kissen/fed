package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.cs.fau.de/kissen/fed/db"
	"io/ioutil"
	"net/http"
	"strconv"
)

type FedPost = db.FedPost
type FedUser = db.FedUser

func makeHandleRootGet(db db.FedStorer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func makeHandleUserPostGet(db db.FedStorer) func(http.ResponseWriter, *http.Request) {
	getArgs := func(r *http.Request) (uint64, uint64, error) {
		vars := mux.Vars(r)

		userId, err := strconv.ParseUint(vars["userId"], 10, 64)
		if err != nil {
			return 0, 0, err
		}

		postId, err := strconv.ParseUint(vars["postId"], 10, 64)
		if err != nil {
			return userId, 0, err
		}

		return userId, postId, nil
	}

	getPost := func(r *http.Request) (*FedPost, error) {
		userId, postId, err := getArgs(r)
		if err != nil {
			return nil, err
		}

		user := db.GetUser(userId)
		post := db.GetPost(postId)

		if user == nil || post == nil {
			err := errors.New("not found")
			return post, err
		}

		if user.Id != post.UserId {
			err := errors.New("user did not author that post")
			return post, err
		}

		return post, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		post, err := getPost(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		json, _ := json.Marshal(post)
		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func makeHandleUserPostPut(db db.FedStorer) func(http.ResponseWriter, *http.Request) {
	getUserFromUrl := func(r *http.Request) (*FedUser, error) {
		vars := mux.Vars(r)

		userId, err := strconv.ParseUint(vars["userId"], 10, 64)
		if err != nil {
			return nil, err
		}

		user := db.GetUser(userId)

		if user == nil {
			err := errors.New("user not found")
			return nil, err
		}

		return user, nil
	}

	getPostFromJson := func(r *http.Request) (*FedPost, error) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		post := &FedPost{}
		err = json.Unmarshal(body, post)
		if err != nil {
			return nil, err
		}

		return post, nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		urlUser, err := getUserFromUrl(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		jsonPost, err := getPostFromJson(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if urlUser.Id != jsonPost.UserId {
			err := errors.New("mismatching userId in URL and JSON")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.AddPost(jsonPost.UserId, jsonPost.Content)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Created")
	}
}
