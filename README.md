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

## Repository

The main repository for this code is at [https://codeberg.org/pdwerryhouse/terraform-provider-immich](codeberg).

The github repository is only to allow access to the Hashicorp registry.

## Introduction

This is a Terraform provider for [https://immich.app/](Immich), the self-hosted photo and video management system. 

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

## Example

The following example sets all of the main configuration settings for an Immich server

```terraform
terraform {
  required_providers {
    immich = {
      source  = "pdwerryhouse/immich"
      version = "~> 0.0.1"
    }
  }
}

provider "immich" {
  endpoint = "http://localhost:2283/api"
  apikey   = "some_api_key"
}

resource "immich_config" "config" {
  backups = {
    database = {
      cron_expression = "0 02 * * *"

      enabled = true

      keep_last_amount = 14

    }
  }
  ffmpeg = {
    accel = "disabled"

    accel_decode = true

    accepted_audio_codecs = ["aac", "mp3", "opus"]

    accepted_containers = ["mov", "ogg", "webm"]

    accepted_video_codecs = ["h264"]

    bframes = -1

    cq_mode = "auto"

    crf = 23

    gop_size = 0

    max_bitrate = "0"

    preferred_hw_device = "auto"

    preset = "ultrafast"

    realtime = {
      enabled = true

      resolutions = [480, 720, 1080]

      video_codecs = ["h264", "hevc"]

    }
    refs = 0

    target_audio_codec = "aac"

    target_resolution = "720"

    target_video_codec = "h264"

    temporal_aq = false

    threads = 0

    tonemap = "hable"

    transcode = "required"

    two_pass = false

  }
  image = {
    colorspace = "p3"

    extract_embedded = false

    fullsize = {
      enabled = false

      format = "jpeg"

      progressive = false

      quality = 80

    }
    preview = {
      format = "jpeg"

      progressive = false

      quality = 80

      size = 1440

    }
    thumbnail = {
      format = "webp"

      progressive = false

      quality = 80

      size = 250

    }
  }
  integrity_checks = {
    checksum_files = {
      cron_expression = "0 03 * * *"

      enabled = true

      percentage_limit = 1

      time_limit = 3600000

    }
    missing_files = {
      cron_expression = "0 03 * * *"

      enabled = true

    }
    untracked_files = {
      cron_expression = "0 03 * * *"

      enabled = true

    }
  }
  job = {
    background_task = {
      concurrency = 5

    }
    editor = {
      concurrency = 2

    }
    face_detection = {
      concurrency = 2

    }
    integrity_check = {
      concurrency = 1

    }
    library = {
      concurrency = 5

    }
    metadata_extraction = {
      concurrency = 5

    }
    migration = {
      concurrency = 5

    }
    notifications = {
      concurrency = 5

    }
    ocr = {
      concurrency = 1

    }
    search = {
      concurrency = 5

    }
    sidecar = {
      concurrency = 5

    }
    smart_search = {
      concurrency = 2

    }
    thumbnail_generation = {
      concurrency = 3

    }
    video_conversion = {
      concurrency = 1

    }
    workflow = {
      concurrency = 5

    }
  }
  library = {
    scan = {
      cron_expression = "0 0 * * *"

      enabled = true

    }
    watch = {
      enabled = false

    }
  }
  logging = {
    enabled = true

    level = "log"

  }
  machine_learning = {
    availability_checks = {
      enabled = true

      interval = 30000

      timeout = 2000

    }
    clip = {
      enabled = true

      model_name = "ViT-B-32__openai"

    }
    duplicate_detection = {
      enabled = true

      max_distance = 0.01

    }
    enabled = true

    facial_recognition = {
      enabled = true

      max_distance = 0.5

      min_faces = 3

      min_score = 0.7

      model_name = "buffalo_l"

    }
    ocr = {
      enabled = true

      max_resolution = 736

      min_detection_score = 0.5

      min_recognition_score = 0.8

      model_name = "PP-OCRv5_mobile"

    }
    urls = ["http://immich-machine-learning:3003"]

  }
  map = {
    dark_style = "https://tiles.immich.cloud/v1/style/dark.json"

    enabled = true

    light_style = "https://tiles.immich.cloud/v1/style/light.json"

  }
  metadata = {
    faces = {
      import = false

    }
  }
  new_version_check = {
    channel = "stable"

    enabled = true

  }
  nightly_tasks = {
    cluster_new_faces = true

    database_cleanup = true

    generate_memories = true

    missing_thumbnails = true

    start_time = "00:00"

    sync_quota_usage = true

  }
  notifications = {
    smtp = {
      enabled = false

      from = ""

      reply_to = ""

      transport = {
        host = ""

        ignore_cert = false

        password = ""

        port = 587

        secure = false

        username = ""

      }
    }
  }
  oauth = {
    account_management_url = ""

    allow_insecure_requests = false

    auto_launch = false

    auto_register = true

    button_text = "Login with OAuth"

    client_id = ""

    client_secret = ""

    default_storage_quota = 0

    enabled = false

    end_session_endpoint = ""

    issuer_url = ""

    mobile_override_enabled = false

    mobile_redirect_uri = ""

    profile_signing_algorithm = "none"

    prompt = ""

    role_claim = "immich_role"

    scope = "openid email profile"

    signing_algorithm = "RS256"

    storage_label_claim = "preferred_username"

    storage_quota_claim = "immich_quota"

    timeout = 30000

    token_endpoint_auth_method = "client_secret_post"

  }
  password_login = {
    enabled = true

  }
  reverse_geocoding = {
    enabled = true

  }
  server = {
    external_domain = ""

    login_page_message = ""

    public_users = true

  }
  storage_template = {
    enabled = false

    hash_verification_enabled = true

    template = "{{y}}/{{y}}-{{MM}}-{{dd}}/{{filename}}"

  }
  templates = {
    email = {
      album_invite_template = ""

      album_update_template = ""

      welcome_template = ""

    }
  }
  theme = {
    custom_css = ""

  }
  trash = {
    days = 30

    enabled = true

  }
  user = {
    delete_delay = 8

  }
}
```

## Contributing code

LLM contributions are not welcome.

## Bug reports

LLM generated bug reports are not welcome.

## Licence

GNU GENERAL PUBLIC [LICENSE](LICENSE) V3 or later