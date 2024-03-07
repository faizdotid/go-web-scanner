package app

type Exploit struct {
	Type               string `json:"type"`
	Description        string `json:"description"`
	FilePath           string `json:"file_path"`
	ValidationCriteria string `json:"validation_criteria"`
	SaveAs             string `json:"save_as"`
	Response           string `json:"response"`
}

type Configuration struct {
	Timeout        int      `json:"timeout"`
	RequestHeaders []string `json:"request_headers"`
}

type Config struct {
	Exploits      []Exploit     `json:"exploits"`
	Configuration Configuration `json:"configuration"`
}
