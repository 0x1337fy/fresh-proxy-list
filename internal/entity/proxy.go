package entity

type Proxy struct {
	Category  string  `json:"category"`
	Proxy     string  `json:"proxy"`
	IP        string  `json:"ip"`
	Port      string  `json:"port"`
	TimeTaken float64 `json:"time_taken"`
	CheckedAt string  `json:"checked_at"`
}

type AdvancedProxy struct {
	Proxy      string   `json:"proxy"`
	IP         string   `json:"ip"`
	Port       string   `json:"port"`
	TimeTaken  float64  `json:"time_taken"`
	CheckedAt  string   `json:"checked_at"`
	Categories []string `json:"categories"`
}
