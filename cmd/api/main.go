package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"net/url"

	"time"

	"github.com/alioygur/cloudinary-go"
	"github.com/alioygur/fb-tinder-app/api"
	"github.com/alioygur/fb-tinder-app/providers"
	fbrepo "github.com/alioygur/fb-tinder-app/providers/fb"
	"github.com/alioygur/fb-tinder-app/providers/fbmock"
	"github.com/alioygur/fb-tinder-app/providers/mongo"
	mysqlrepo "github.com/alioygur/fb-tinder-app/providers/mysql"
	services "github.com/alioygur/fb-tinder-app/service"
	"github.com/alioygur/goutil"
	"gopkg.in/mgo.v2"
)

func main() {
	if err := checkRequiredEnv(); err != nil {
		log.Fatal(err)
	}

	// waitForServices()

	// db, err := mysqlrepo.ConnectToDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	mongoSess, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoSess.Close()

	// repos
	var fb services.FacebookRepository
	var sql services.Repository

	// sql = mysqlrepo.New(db)
	sql = mongo.New(mongoSess)

	switch goutil.EnvMustGet("APP_ENV") {
	default:
		// instances number of users
		users := mysqlrepo.GenUsers(10)
		_fb := fbmock.New(uint64(goutil.EnvMustInt("FB_APP_ID")), goutil.EnvMustGet("FB_APP_SECRET"))
		// seed facebook mock repo with fake users
		_fb.Users = users
		fb = _fb
	case "production":
		fb = fbrepo.New(uint64(goutil.EnvMustInt("FB_APP_ID")), goutil.EnvMustGet("FB_APP_SECRET"))
	}

	// deps
	jwt := providers.NewJWT()
	cc, err := cloudinary.New(goutil.EnvMustGet("CLOUDINARY_URL"))
	if err != nil {
		log.Fatal(err)
	}
	imgCDN := providers.NewImageCDN(cc)

	// services
	service := services.New(fb, sql, jwt, imgCDN)

	h := api.NewHandler(service)
	log.Printf("server starting port: %s env: %s", os.Getenv("PORT"), os.Getenv("APP_ENV"))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), h); err != nil {
		log.Fatal(err)
	}
}

// checks required env variable must be setted
func checkRequiredEnv() error {
	envs := [...]string{
		"APP_ENV",
		"PORT",
		"SECRET_KEY",
		"MYSQL_URL",
		"PROXY_URL",
		"CLOUDINARY_URL",
		"FB_APP_ID",
		"FB_APP_SECRET",
		"FB_REQUIRED_PERMS",
	}

	for _, v := range envs {
		if os.Getenv(v) == "" {
			return fmt.Errorf("the %s env variable required", v)
		}
	}
	return nil
}

func waitForServices() {
	mysql, err := url.Parse(goutil.EnvMustGet("MYSQL_URL"))
	if err != nil {
		log.Fatal(err)
	}

	services := []url.URL{
		url.URL{Scheme: "tcp", Host: mysql.Host},
	}

	if err := goutil.WaitForServices(services, 15*time.Second); err != nil {
		log.Fatal(err)
	}
	log.Println("services are ready")
}
