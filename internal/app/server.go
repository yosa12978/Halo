package app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/yosa12978/halo/internal/pkg/mongo"
)

func Run() {
	if err := godotenv.Load("./config/.env"); err != nil {
		panic(err)
	}

	addr := os.Getenv(fmt.Sprintf("MONGO_ADDR_%s", os.Getenv("MODE")))
	mongo.InitMongo(addr, os.Getenv("DB_NAME"))

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
	}()

	out := make(chan os.Signal, 1)
	signal.Notify(out, syscall.SIGINT, syscall.SIGTERM)
	<-out
	log.Printf("Program stopped at %d\n", time.Now().Unix())
}
