
resource "immich_album" "test1" {
  album_name  = "Test"
  description = "Test Album"
  order       = "desc"
}

resource "immich_album" "test2" {
  album_name  = "Test"
  description = "Test Album"
  order       = "asc"
}

