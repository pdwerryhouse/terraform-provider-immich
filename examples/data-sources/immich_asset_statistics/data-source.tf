
data "immich_asset_statistics" "all" {
}

data "immich_asset_statistics" "favourites" {
  is_favourite = true
}

data "immich_asset_statistics" "hidden" {
  is_hidden = true
}

data "immich_asset_statistics" "archive" {
  visibility = "archive"
}

data "immich_asset_statistics" "timeline" {
  visibility = "timeline"
}

data "immich_asset_statistics" "hidden" {
  visibility = "hidden"
}