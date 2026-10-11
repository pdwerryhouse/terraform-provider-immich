
resource "immich_asset" "image" {
  asset_data  = filebase64("picture.jpg")
  filename    = "image.jpg"
  created_at  = "2026-10-11T12:12:11+11:00"
  modified_at = "2026-10-11T12:12:11+11:00"
}

