package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Register struct {
	gorm.Model
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	UserName  string `json:"username"`
	Passworda string `json:"password"`
}
type LogError struct {
	LError string
}

var tpl *template.Template
var db *gorm.DB
var err error

func main() {
	connectDB()
	initlizeMigrate()
	log.Println("Server running on port 3000")
	handleFunc()
}

func connectDB() {
	db, err = gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

}
func init() {
	tpl = template.Must(template.ParseGlob("web/view/*.html"))
}
func initlizeMigrate() {
	db.AutoMigrate(&Register{})

}

func Index(response http.ResponseWriter, request *http.Request) {

	temp, err := template.ParseFiles("web/view/signup.html")
	if err != nil {
		fmt.Println(err.Error())
		panic("File Open Error")
	}
	temp.Execute(response, nil)
}
func home(response http.ResponseWriter, request *http.Request) {
	fmt.Fprint(response, "Welcome in home Page")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func welcome(response http.ResponseWriter, request *http.Request) {
	temp, err := template.ParseFiles("web/view/signup.html")
	if err != nil {
		fmt.Println(err.Error())
		panic("File Open Error")
	}
	request.ParseForm()
	var firstname string = request.Form.Get("fname")
	lastname := request.Form.Get("lname")
	useremail := request.Form.Get("uemail")
	username := request.Form.Get("uname")
	pwd := request.Form.Get("pwd")
	fmt.Println("First Name ", firstname)
	fmt.Println("Last Name ", lastname)
	fmt.Println("Email ", useremail)
	fmt.Println("Username name ", username)
	fmt.Println("Password ", pwd)
	secret, err := HashPassword(pwd)
	if err != nil {
		return
	}
	fmt.Println("Password ", secret)
	rerr := LogError{
		LError: "Username already exists",
	}
	var allreister Register
	db.Table("registers").Select("first_name", "user_name").Where("user_name = ?", username).Scan(&allreister)
	runame := allreister.UserName
	if runame == username {
		temp.ExecuteTemplate(response, "signup.html", rerr)
	} else {
		db.Create(&Register{FirstName: firstname, LastName: lastname, Email: useremail, UserName: username, Passworda: secret})
		tpl.ExecuteTemplate(response, "hello.html", nil)
	}

}
func allUsers(response http.ResponseWriter, r *http.Request) {
	var allreister []Register
	db.Find(&allreister)
	fmt.Println(allreister)
	json.NewEncoder(response).Encode(allreister)
}
func specificUser(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	vars := mux.Vars(request)

	id := vars["id"]
	//a, _ := strconv.Atoi(id)
	var allreister Register
	lpwd := request.Form.Get("password")
	db.Table("registers").Select("first_name", "last_name", "user_name", "passworda").Where("id = ?", id).Scan(&allreister)
	fmt.Println(allreister)
	json.NewEncoder(response).Encode(allreister)
	fmt.Fprintf(response, allreister.FirstName)
	fmt.Fprintf(response, allreister.UserName)
	fmt.Fprintf(response, allreister.Passworda)
	match := CheckPasswordHash(lpwd, allreister.Passworda)
	fmt.Println("Match ", match)
}

func greet(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	lusername := request.Form.Get("username")
	lpwd := request.Form.Get("password")
	fmt.Println(" username ", lusername+" ", lpwd)

	var allreister Register
	db.Table("registers").Select("first_name", "last_name", "user_name", "passworda").Where("user_name = ?", lusername).Scan(&allreister)
	//json.NewEncoder(response).Encode(allreister)
	msg := LogError{
		LError: "Invalid username and password",
	}
	str := Register{
		FirstName: allreister.FirstName,
		LastName:  allreister.LastName,
	}
	match := CheckPasswordHash(lpwd, allreister.Passworda)
	if match == true {
		tpl.ExecuteTemplate(response, "greeting.html", str)
	} else {
		tpl.ExecuteTemplate(response, "hello.html", msg)
	}

}
func login(response http.ResponseWriter, request *http.Request) {
	tpl.ExecuteTemplate(response, "hello.html", nil)
}

func handleFunc() {
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/signup", Index)
	myRouter.HandleFunc("/signin", login)
	myRouter.HandleFunc("/register", welcome)
	myRouter.HandleFunc("/welcome", greet)
	myRouter.HandleFunc("/users", allUsers)
	myRouter.HandleFunc("/user/{id}", specificUser)

	http.ListenAndServe(":3000", myRouter)
}
