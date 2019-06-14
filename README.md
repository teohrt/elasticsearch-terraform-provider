# terraform-provider-esslowlogconfig
This is a terraform provider that lets you provision elasticsearch slow log configuration for indexes at the shard level.

## Installation

Build the binary, and put it in a good spot on your system. Then update your `~/.terraformrc` to refer to the binary:

```hcl
providers {
  esslowlogconfig = "/path/to/terraform-provider-esslowlogconfig"
}
```

See [the docs for more information](https://www.terraform.io/docs/plugins/basics.html).

## Usage

```tf
provider "esslowlogconfig" {
    es_endpoint = "https://search-foo-bar-pqrhr4w3u4dzervg41frow4mmy.us-east-1.es.amazonaws.com" # Don't include port at the end for aws
    region = "us-east-2"
}

resource "esslowlogconfig_index" "index_1" {
    name = "activity_points"

    query_warn_threshold = "10s"
    query_info_threshold = "5s"
    query_debug_threshold = "2s"
    query_trace_threshold = "500ms"

    fetch_warn_threshold = "1s"
    fetch_info_threshold = "800ms"
    fetch_debug_threshold = "500ms"
    fetch_trace_threshold = "200ms"
}

resource "esslowlogconfig_index" "index_2" {
    name = ".kibana_1"

    query_warn_threshold = "10s"
    query_info_threshold = "5s"
    query_debug_threshold = "2s"
    query_trace_threshold = "500ms"

    fetch_warn_threshold = "1s"
    fetch_info_threshold = "800ms"
    fetch_debug_threshold = "500ms"
    fetch_trace_threshold = "200ms"
}
```

#### Environment variables

You can provide your credentials via the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`, environment variables, representing your AWS Access Key and AWS Secret Key. If applicable, the `AWS_SESSION_TOKEN` environment variables is also supported.

Example usage:

```shell
$ export AWS_ACCESS_KEY_ID="anaccesskey"
$ export AWS_SECRET_ACCESS_KEY="asecretkey"
$ export AWS_SESSION_TOKEN="asessiontoken"
$ terraform plan
```

## Development

### Requirements

* [Golang](https://golang.org/dl/) >= 1.11


```
go build -o /path/to/binary/terraform-provider-esslowlogconfig
```

## Licence

See LICENSE.

## Contributing

1. Fork it ( https://github.com/phillbaker/terraform-provider-elasticsearch/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request