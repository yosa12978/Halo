package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/yosa12978/halo/internal/api"
	"github.com/yosa12978/halo/internal/pkg/mongo"
)

func Run() {
	if err := godotenv.Load("./config/.env"); err != nil {
		panic(err)
	}

	addr := os.Getenv(fmt.Sprintf("MONGO_ADDR_%s", os.Getenv("MODE")))
	mongo.InitMongo(addr, os.Getenv("DB_NAME"))

	api.Run()

	out := make(chan os.Signal, 1)
	signal.Notify(out, syscall.SIGINT, syscall.SIGTERM)
	sig := <-out
	log.Printf("Program stopped at %d signal: %s\n", time.Now().Unix(), sig.String())
}
