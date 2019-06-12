provider "esdynamiconfig" {
    es_endpoint = "http://localhost:8080/"
    region = "us-east-2"
}

resource "esdynamiconfig_index" "test" {
    indexName = "test"
    query_warn_threshold = "5s"
    query_info_threshold = "2s"
}