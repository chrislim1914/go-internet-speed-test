package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chrislim1914/go-internet-speed-test/speedtester"
	"github.com/chrislim1914/go-internet-speed-test/utilities"
)

func main() {
	apiMode := flag.Bool("api", false, "Run as HTTP API server")
	port := flag.String("port", "8080", "API server port")
	flag.Parse()

	if *apiMode {
		// TODO: implement REST and WebSocket API server
		// REST API example endpoint
		// GetBestServer and GetServerList
		// WebSocket for real-time speed test updates
		// TestDownloadSpeed and TestUploadSpeed
		http.HandleFunc("/speedtest", func(w http.ResponseWriter, r *http.Request) {
			provider := &speedtester.DefaultSpeedTestProvider{}
			serverList, err := provider.GetServerList()
			if err != nil || len(serverList) == 0 {
				http.Error(w, "No servers available", http.StatusInternalServerError)
				return
			}
			bestServer, err := provider.GetBestServer(serverList)
			if err != nil {
				http.Error(w, "No best server found", http.StatusInternalServerError)
				return
			}
			downloadSpeeds, _, err := speedtester.TestDownloadSpeed(*bestServer)
			if err != nil {
				http.Error(w, "Download test failed", http.StatusInternalServerError)
				return
			}
			uploadSpeeds, _, err := speedtester.TestUploadSpeed(*bestServer)
			if err != nil {
				http.Error(w, "Upload test failed", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Download: %v Mbps\nUpload: %v Mbps\n", downloadSpeeds, uploadSpeeds)
		})
		log.Printf("API server running on :%s", *port)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
		os.Exit(0)
	}

	// CLI mode (your existing code)
	provider := &speedtester.DefaultSpeedTestProvider{}
	stop := utilities.StartSpinner("Fetching server list...")
	serverList, err := provider.GetServerList()
	stop()
	if err != nil {
		panic(err)
	}
	stop = utilities.StartSpinner("Finding best server...")
	bestServer, err := provider.GetBestServer(serverList)
	stop()
	if err != nil {
		panic(err)
	}
	if bestServer != nil {
		println("Best server:", bestServer.Name, "with URL:", bestServer.Server)
	} else {
		println("No reachable server found.")
	}
	stop = utilities.StartSpinner("Fetching ISP information...")
	ispInfo, err := provider.GetISPInfo()
	stop()
	if err != nil {
		panic(err)
	}
	fmt.Printf("ISP Provider: %+v\n", ispInfo.Organization)
	stop = utilities.StartSpinner("Testing download speed...")
	downloadSpeeds, _, err := speedtester.TestDownloadSpeed(*bestServer)
	stop()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Download Speeds (Mbps): %v\n", downloadSpeeds)
	stop = utilities.StartSpinner("Testing upload speed...")
	uploadSpeeds, _, err := speedtester.TestUploadSpeed(*bestServer)
	stop()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Upload Speeds (Mbps): %v\n", uploadSpeeds)
}
