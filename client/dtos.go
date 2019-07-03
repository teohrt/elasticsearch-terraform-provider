package client

type GetItemResponse struct {
	Query_warn_threshold  string
	Query_info_threshold  string
	Query_debug_threshold string
	Query_trace_threshold string

	Fetch_warn_threshold  string
	Fetch_info_threshold  string
	Fetch_debug_threshold string
	Fetch_trace_threshold string
}

type Container struct {
	Settings Settings `json:"settings"`
}

type Settings struct {
	Index Index `json:"index"`
}

type Index struct {
	Search Search `json:"search"`
}

type Search struct {
	Slowlog Slowlog `json:"slowlog"`
}

type Slowlog struct {
	Threshold Threshold `json:"threshold"`
}

type Threshold struct {
	Query ThresholdType `json:"query"`
	Fetch ThresholdType `json:"fetch"`
}

type ThresholdType struct {
	Warn  string `json:"warn"`
	Info  string `json:"info"`
	Debug string `json:"debug"`
	Trace string `json:"trace"`
}
