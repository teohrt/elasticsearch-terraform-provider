provider "esdynamiconfig" {
    es_endpoint = "https://search-activitygraph-teo-rfbftdx5r3c5rwqoxjkrphrriy.us-east-2.es.amazonaws.com"
    region = "us-east-2"
}

resource "esdynamiconfig_index" "test" {
    name = "activity_points"
    query_warn_threshold = "100s"
    query_info_threshold = "100s"
}