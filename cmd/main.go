package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"ozon/pkg/handlers"
	"ozon/pkg/middleware"

	// "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	//if err := initConfig(); err != nil {
	//	fmt.Printf("error initializing config : %s\n", err.Error())
	//	return
	//}

	port := "8080"
	storage := "postgres"

	zapLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("can't zap.NewProduction()")
		return
	}
	defer func() {
		err = zapLogger.Sync()
		fmt.Println("can't  zapLogger.Sync()")
	}()

	r := mux.NewRouter()

	var db *sql.DB
	//if viper.Get("storage").(string) == "postgres" {
	if storage == "postgres" {
		dsn := "postgres://postgres:root@localhost:5432/ozon_test?sslmode=disable"

		db, err = sql.Open("postgres", dsn)
		if err != nil {
			fmt.Println("errpr!!!!", err)
			return
		}

		db.SetMaxOpenConns(10)

		err = db.Ping()
		if err != nil {
			fmt.Println("no base")
			panic(err)
		}
	} else {
		db = nil
	}

	LinkHandler := handlers.NewLinkHandler(db, storage)

	r.HandleFunc("/api/newLink", LinkHandler.NewLink).Methods("POST")
	r.HandleFunc("/api/getLink", LinkHandler.GetLink).Methods("GET")

	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(".././static"))))

	handler := middleware.LofInfo(r)
	//port := viper.Get("port")

	fmt.Println("start serv on port " + port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Println("can't Listen and server")
		return
	}
}

//func initConfig() error {
//	viper.AddConfigPath("../configs")
//	viper.SetConfigName("config")
//	return viper.ReadInConfig()
//}
