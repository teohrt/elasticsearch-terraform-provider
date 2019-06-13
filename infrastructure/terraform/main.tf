provider "esdynamiconfig" {
    es_endpoint = "https://search-activitygraph-teo-rfbftdx5r3c5rwqoxjkrphrriy.us-east-2.es.amazonaws.com"
    region = "us-east-2"
}

resource "esdynamiconfig_index" "index_1" {
    name = "activity_points"
    query_warn_threshold = "5s"
    query_info_threshold = "2s"
}

resource "esdynamiconfig_index" "index_2" {
    name = ".kibana_1"
    query_warn_threshold = "10s"
    query_info_threshold = "5s"
}