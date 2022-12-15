package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/souravdev-eng/mongocrud/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

const DBName = "mogo-golang"

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userId := bson.ObjectIdHex(id)

	user := models.User{}
	if err := uc.session.DB(DBName).C("users").FindId(userId).One(&user); err != nil {
		w.WriteHeader(404)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := models.User{}

	json.NewDecoder(r.Body).Decode(&user)

	user.Id = bson.NewObjectId()
	uc.session.DB(DBName).C("users").Insert(user)
	uj, err := json.Marshal(user)

	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", uj)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	userId := bson.ObjectIdHex(id)

	if err := uc.session.DB(DBName).C("users").RemoveId(userId); err != nil {
		w.WriteHeader(404)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Delete user", userId, "\n")
}
