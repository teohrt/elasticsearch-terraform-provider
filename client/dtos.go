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
