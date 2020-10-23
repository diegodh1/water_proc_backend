package api

import (
	"fmt"
	"log"
	"net/http"

	conf "water_proccesing/config"
	"water_proccesing/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//App struct
type App struct {
	Router *mux.Router
	DB     *gorm.DB
	Config *conf.Config
}

//Initialize initialize db
func (a *App) Initialize(config *conf.Config) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Database,
	)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	fmt.Println(dsn)
	if err != nil {
		log.Fatal("Could not connect database")
	}
	log.Println("Connected!")
	a.DB = db
	a.Config = config
	a.Router = mux.NewRouter()
	a.setRouters()
}

//setting the routes
func (a *App) setRouters() {
	a.Post("/userLogin", a.UserLogin)
	a.Post("/searchUser", a.SearchUser)
	a.Get("/getCompanies", a.GetCompanies)
	a.Post("/createUser", a.CreateUser)
	a.Put("/updateUser/{UserID}", a.UpdateUser)
	a.Post("/createProfile", a.CreateProfile)
	a.Put("/updateProfile/{ProfileID}", a.UpdateProfile)
	a.Post("/createUserProfile", a.CreateUserProfile)
	a.Post("/updateUserProfile", a.UpdateUserProfile)
}

//Get all get functions
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

//Post all Post functions
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

//Put all Put functions
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

//UserLogin let user login
func (a *App) UserLogin(w http.ResponseWriter, r *http.Request) {
	handler.UserLogin(a.DB, w, r, a.Config.SecretSeed)
}

//GetCompanies get all companies at the database
func (a *App) GetCompanies(w http.ResponseWriter, r *http.Request) {
	handler.GetCompanies(a.DB, w, r)
}

//SearchUser search a user by ID
func (a *App) SearchUser(w http.ResponseWriter, r *http.Request) {
	handler.SarchUser(a.DB, w, r)
}

//CreateUser create a new user
func (a *App) CreateUser(w http.ResponseWriter, r *http.Request) {
	handler.CreateUser(a.DB, w, r)
}

//UpdateUser let update a user
func (a *App) UpdateUser(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUser(a.DB, w, r)
}

//CreateProfile create a new profile
func (a *App) CreateProfile(w http.ResponseWriter, r *http.Request) {
	handler.CreateProfile(a.DB, w, r)
}

//UpdateProfile let update a user
func (a *App) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	handler.UpdateProfile(a.DB, w, r)
}

//CreateUserProfile let assing a profile to a user
func (a *App) CreateUserProfile(w http.ResponseWriter, r *http.Request) {
	handler.CreateUserProfile(a.DB, w, r)
}

//UpdateUserProfile let update a profile given to a user
func (a *App) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	handler.UpdateUserProfile(a.DB, w, r)
}

//Run run app
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(a.Router)))
}
