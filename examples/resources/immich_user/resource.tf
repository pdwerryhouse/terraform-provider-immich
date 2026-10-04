resource "immich_user" "user" {
  name                = "john"
  email               = "john@example.com"
  password            = "some_secret_password"
  is_admin            = false
  notify              = false
  quota_size_in_bytes = 20000000
  avatar_color        = "yellow"
}

