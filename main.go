package main

import (
	"fmt"

	"github.com/Maltide/link-shortener/pkg/config"
	"github.com/Maltide/link-shortener/pkg/logger"
)

// import (
// 	"fmt"
// 	"net/http"
// 	"time"
// )

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Custom server configuration!")
// }

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Errorf("fail to create config:", err)
		return
	}
	log, err := logger.Logger("debug")
	if err != nil {
		fmt.Errorf("fail to create logger: ", err)
		return
	}
	// server := &http.Server{
	// 	Addr:         ":8080",
	// 	Handler:      http.HandlerFunc(handler),
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }

	// fmt.Println("Starting custom server at port 8080")
	// err := server.ListenAndServe()
	// if err != nil {
	// 	fmt.Println("Error starting the server:", err)
	// }
}
