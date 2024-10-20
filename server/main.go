package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/tamakoshi2001/gextension/handler/router"
	"github.com/tamakoshi2001/gextension/model"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
		sitesPath     = "save/sites.json"
		vectorsPath   = "save/vectors.json"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	sites, err := loadSites(sitesPath)
	if err != nil {
		log.Println("Error reading file: ", err)
		sites = []model.Site{}
	}
	vectors, err := loadVectors(vectorsPath)
	if err != nil {
		log.Println("Error reading file: ", err)
		vectors = []model.Vector{}
	}
	if len(sites) != len(vectors) {
		sites = []model.Site{}
		vectors = []model.Vector{}
	}
	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(&sites, &vectors)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    port,
		Handler: cors.AllowAll().Handler(mux),
	}

	var wg sync.WaitGroup

	// Start the server in a goroutine so that it doesn't block.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// It's important to check for http.ErrServerClosed to not catch a server
			// intentionally closed.
			log.Fatalf("ListenAndServe(): %s", err)
		}
		// time.Sleep(10 * time.Second)
		log.Println("goroutine: server is shutdown")
	}()

	// チャネルに入ってきた情報によって処理をわける
	<-ctx.Done()

	log.Println("main: timeout context is started")
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown(): %s", err)
	}
	log.Println("main: server is shutdown")

	// Wait for the goroutine (not only server) to shutdown.
	wg.Wait()

	err = saveSites(sites, sitesPath)
	if err != nil {
		return err
	}
	err = saveVectors(vectors, vectorsPath)
	if err != nil {
		return err
	}

	// print shutdown message
	log.Println("main: server shutdown successfully")

	return nil
}

// save sites  to a file
func saveSites(sites []model.Site, sitesPath string) error {

	// 構造体をJSONに変換
	data, err := json.MarshalIndent(sites, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON: ", err)
		return nil
	}

	// JSONをファイルに保存
	err = os.WriteFile(sitesPath, data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file: ", err)
		return err
	}
	return nil
}

// save vectors  to a file
func saveVectors(vectors []model.Vector, vectorsPath string) error {

	// 構造体をJSONに変換
	data, err := json.MarshalIndent(vectors, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON: ", err)
		return nil
	}

	// JSONをファイルに保存
	err = os.WriteFile(vectorsPath, data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file: ", err)
		return err
	}
	return nil
}

// load sites from a file
func loadSites(sitesPath string) ([]model.Site, error) {
	// ファイルを読み込む
	data, err := os.ReadFile(sitesPath)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return nil, err
	}

	// JSONを構造体に変換
	var sites []model.Site
	err = json.Unmarshal(data, &sites)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return nil, err
	}

	return sites, nil
}

// load vectors from a file
func loadVectors(vectorsPath string) ([]model.Vector, error) {
	// ファイルを読み込む
	data, err := os.ReadFile(vectorsPath)
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return nil, err
	}

	// JSONを構造体に変換
	var vectors []model.Vector
	err = json.Unmarshal(data, &vectors)
	if err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return nil, err
	}

	return vectors, nil
}
