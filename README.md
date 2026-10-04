<!--
 Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
 
 This file is part of terraform-provider-immich.
 
 terraform-provider-immich is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.
 
 terraform-provider-immich is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.
 
 You should have received a copy of the GNU General Public License
 along with terraform-provider-immich.  If not, see <https://www.gnu.org/licenses/>.
-->

# Terraform Provider for Immich

## Installation

```terraform
terraform {
  required_providers {
    immich = {
      source  = "pdwerryhouse/immich"
      version = "~> 0.0.1"
    }
  }
}
```

## Usage

```terraform

provider "immich" {
  endpoint = "http://localhost:2283/api"
  apikey   = "some_api_key"
}

resource "immich_album" "test" {
  album_name  = "Test"
  description = "Test Album"
}
```

## Testing

To run tests, set the IMMICH_ENDPOINT and IMMICH_API_KEY environment variables:

```bash
export IMMICH_ENDPOINT=http://localhost:2283/api
read -s IMMICH_API_KEY
export IMMICH_API_KEY
```

Then run the following:

```bash
TF_ACC=1 go test -count=1 -v ./...
```

## Licence

GNU GENERAL PUBLIC [LICENSE](LICENSE) V3 or later