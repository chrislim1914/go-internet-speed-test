package speedtester

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/chrislim1914/go-internet-speed-test/utilities"
)

type Server struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Server      string `json:"server"`
	DownloadURL string `json:"dlURL"`
	UploadURL   string `json:"ulURL"`
	PingURL     string `json:"pingURL"`
	GetIPURL    string `json:"getIpURL"`
	SponsorName string `json:"sponsorName"`
	SponsorURL  string `json:"sponsorURL"`
}

type ISPInfo struct {
	IPAddress    string    `json:"ip_address"`
	ISP          string    `json:"isp"`
	ASN          string    `json:"asn"`
	Organization string    `json:"organization"`
	Country      string    `json:"country"`
	City         string    `json:"city"`
	Region       string    `json:"region"`
	Timestamp    time.Time `json:"timestamp"`
}

type SpeedTestProvider interface {
	GetBestServer([]Server) (*Server, error)
	GetServerList() ([]Server, error)
	GetISPInfo() (*ISPInfo, error)
	TestDownloadSpeed(Server) ([]float64, []float64, error)
	TestUploadSpeed(Server) ([]float64, []float64, error)
}

type DefaultSpeedTestProvider struct{}

func (d *DefaultSpeedTestProvider) GetBestServer(servers []Server) (*Server, error) {
	type serverResult struct {
		server  Server
		latency time.Duration
	}

	results := make(chan serverResult, len(servers))
	var wg sync.WaitGroup

	for _, srv := range servers {
		wg.Add(1)
		go func(s Server) {
			defer wg.Done()
			hostName := utilities.NormalizedURL(s.Server)
			latency, err := utilities.TestServerLatency(hostName)
			if err != nil {
				results <- serverResult{server: s, latency: 10 * time.Second}
				return
			}
			results <- serverResult{server: s, latency: latency}
		}(srv)
	}

	wg.Wait()
	close(results)

	var ranked []serverResult
	for res := range results {
		ranked = append(ranked, res)
	}

	if len(ranked) == 0 {
		return nil, fmt.Errorf("no working servers found")
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].latency < ranked[j].latency
	})

	return &ranked[0].server, nil
}

func (d *DefaultSpeedTestProvider) GetServerList() ([]Server, error) {
	req, err := http.NewRequest("GET", utilities.ServerListUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var servers []Server
	if err := json.Unmarshal(b, &servers); err != nil {
		return nil, err
	}

	return servers, nil
}

func (d *DefaultSpeedTestProvider) GetISPInfo() (*ISPInfo, error) {
	req, err := http.NewRequest("GET", utilities.GetISPInfoUrl, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse struct {
		IP      string `json:"ip"`
		Org     string `json:"org"`
		ASN     string `json:"asn"`
		Country string `json:"country_name"`
		City    string `json:"city"`
		Region  string `json:"region"`
	}

	if err := json.Unmarshal(b, &apiResponse); err != nil {
		return nil, err
	}

	return &ISPInfo{
		IPAddress:    apiResponse.IP,
		ISP:          apiResponse.Org,
		ASN:          apiResponse.ASN,
		Organization: apiResponse.Org,
		Country:      apiResponse.Country,
		City:         apiResponse.City,
		Region:       apiResponse.Region,
		Timestamp:    time.Now(),
	}, nil
}
