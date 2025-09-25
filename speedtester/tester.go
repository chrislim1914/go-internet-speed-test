package speedtester

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chrislim1914/go-internet-speed-test/utilities"
)

func TestDownloadSpeed(server Server) ([]float64, []float64, error) {
	iterations := utilities.Test_Iterations
	duration := utilities.Test_Duration
	speeds := make([]float64, 0, iterations)
	elapses := make([]float64, 0, iterations)

	for i := 0; i < iterations; i++ {

		speed, elapse, err := testDownloadSpeed(duration, server)
		if err != nil {
			fmt.Printf("\nDownload test %d failed: %v\n", i+1, err)
			continue
		}

		speeds = append(speeds, speed)
		elapses = append(elapses, elapse)
		time.Sleep(1 * time.Second)
	}

	if len(speeds) == 0 {
		return nil, nil, fmt.Errorf("all download tests failed")
	}

	return speeds, elapses, nil
}

func testDownloadSpeed(duration time.Duration, server Server) (float64, float64, error) {
	host := utilities.NormalizedURL(server.Server)

	client := &http.Client{Timeout: utilities.Test_Timeout}
	req, err := http.NewRequest("GET", host+server.DownloadURL, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	timer := time.NewTimer(duration)
	defer timer.Stop()

	totalBytes := int64(0)
	buffer := make([]byte, 32*1024)
	start := time.Now()

downloadLoop:
	for {
		select {
		case <-timer.C:
			break downloadLoop
		default:
			n, err := resp.Body.Read(buffer)
			totalBytes += int64(n)
			if err == io.EOF {
				break downloadLoop
			}
			if err != nil {
				return 0, 0, err
			}
		}
	}

	elapsed := time.Since(start).Seconds()
	speedMbps := (float64(totalBytes) * 8) / (elapsed * 1000 * 1000)
	stop := utilities.StartSpinner(fmt.Sprintf("%.2f Mbps in %.2f seconds", speedMbps, elapsed))
	stop()
	return speedMbps, elapsed, nil
}

func TestUploadSpeed(server Server) ([]float64, []float64, error) {
	iterations := utilities.Test_Iterations
	speeds := make([]float64, 0, iterations)
	elapses := make([]float64, 0, iterations)

	for i := 0; i < iterations; i++ {

		speed, elapse, err := testUploadSpeed(1*1024*1024, server)
		if err != nil {
			fmt.Printf("\nUpload test %d failed: %v\n", i+1, err)
			continue
		}

		speeds = append(speeds, speed)
		elapses = append(elapses, elapse)
		time.Sleep(1 * time.Second)
	}

	if len(speeds) == 0 {
		return nil, nil, fmt.Errorf("all upload tests failed")
	}

	return speeds, elapses, nil
}

func testUploadSpeed(dataSize int, server Server) (float64, float64, error) {
	host := utilities.NormalizedURL(server.Server)

	testData := make([]byte, dataSize)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	client := &http.Client{Timeout: utilities.Test_Timeout}
	start := time.Now()

	req, err := http.NewRequest("POST", host+server.UploadURL, bytes.NewReader(testData))
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start).Seconds()
	speedMbps := (float64(dataSize) * 8) / (elapsed * 1000 * 1000)
	stop := utilities.StartSpinner(fmt.Sprintf("%.2f Mbps in %.2f seconds", speedMbps, elapsed))
	stop()
	return speedMbps, elapsed, nil
}
