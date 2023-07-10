# Terraform Provider perimeter81

This provider is built to handel terraform for perimeter81 [perimeter81](https://www.perimeter81.com/)

## Requirements

Terraform 1.3.x
Go 1.19.1 (to build the provider plugin)

## Installation

```shell
make
```

## Usage

```terraform
provider "perimeter81" {
  api_key = "XXXXXXXXXXXX"
  base_url = "https://api.perimeter81.com/api/rest"
}
```

## Run Tests

Run the following:

```shell
PERIMETER81_API_KEY=XXXXXX BASE_URL=https://api.perimeter81.com/api/rest TF_ACC=1 go test -timeout 120m  -parallel 10
```

Note: you need to replace the xxxxxxxxx with your api key.
