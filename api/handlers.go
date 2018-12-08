package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ASeegull/secrets-vault/models"

	uuid "github.com/satori/go.uuid"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
)

const maxSecretLen = 2048

type errResp struct {
	Message string `json:"message"`
	Meta    string `json:"meta"`
}

func respJSON(msg interface{}, w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Warn("failed to send responce:", err)
	}
}

func (a API) getSecret(w http.ResponseWriter, r *http.Request) {
	hash, err := uuid.FromString(chi.URLParam(r, "hash"))
	if err != nil {
		respJSON(errResp{Message: "Missing secret's hash"}, w, http.StatusBadRequest)
		return
	}

	secret, err := a.storage.Get(hash)
	if err != nil {
		respJSON(errResp{Message: "Wrong secret's hash"}, w, http.StatusBadRequest)
		return
	}

	if secret.ExpireAfterViews == 0 || time.Now().After(secret.ExpAfter) {
		respJSON(errResp{Message: "This secret is no longer available", Meta: hash.String()}, w, http.StatusBadRequest)
		return
	}

	secret.ExpireAfterViews, err = a.storage.DecrementViews(hash)
	if err != nil {
		log.Warn(err)
	}

	respJSON(secret, w, http.StatusOK)
}

func (a API) addSecret(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("message")
	if msg == "" {
		respJSON(errResp{Message: "Please, add secret to store"}, w, http.StatusBadRequest)
		return
	}

	if len(msg) > maxSecretLen {
		respJSON(errResp{Message: "Secret's text too long"}, w, http.StatusBadRequest)
		return
	}

	viewsAllowed, err := strconv.Atoi(r.FormValue("expireAfterViews"))
	if err != nil || viewsAllowed == 0 {
		respJSON(errResp{Message: "Please, specify max views count"}, w, http.StatusBadRequest)
		return
	}

	t, err := strconv.Atoi(r.FormValue("expireAfter"))
	if err != nil || viewsAllowed == 0 {
		respJSON(errResp{Message: "Please, specify how long should we keep your secret"}, w, http.StatusBadRequest)
		return
	}

	expiresAfter := time.Now().Add(time.Minute * time.Duration(t))

	hash, err := uuid.NewV4()
	if err != nil || viewsAllowed == 0 {
		log.Warn("Failed to save secret", "could not create uuid:", err.Error())
		respJSON(errResp{Message: "Failed to save secret"}, w, http.StatusBadRequest)
		return
	}

	if err = a.storage.Save(hash, viewsAllowed, expiresAfter, msg); err != nil {
		log.Warn("Failed to save secret", err.Error())
		respJSON(errResp{Message: "Failed to save secret"}, w, http.StatusBadRequest)
	}

	respJSON(models.Secret{ID: hash, ExpAfter: expiresAfter}, w, http.StatusCreated)
}
