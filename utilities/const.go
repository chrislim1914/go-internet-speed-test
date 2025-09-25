package utilities

import "time"

const (
	ServerListUrl   = `https://librespeed.org/backend-servers/servers.php`
	GetISPInfoUrl   = `https://ipinfo.io/json`
	Test_Iterations = 10
	Test_Duration   = 5 * time.Second
	Test_Timeout    = 30 * time.Second
)
