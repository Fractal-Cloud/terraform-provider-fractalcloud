resource "fractalcloud_personal_bounded_context" "production_audi" {
  short_name   = "production-audi"
  display_name = "Audi Production"
}

resource "fractalcloud_personal_bounded_context" "production_toyota" {
  short_name   = "production-toyota"
  display_name = "Toyota Production"
}

resource "fractalcloud_personal_bounded_context" "production" {
  short_name   = "production"
  display_name = "Production"
}
