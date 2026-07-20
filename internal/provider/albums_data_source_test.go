// Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
//
// This file is part of terraform-provider-immich.
//
// terraform-provider-immich is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// terraform-provider-immich is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with terraform-provider-immich.  If not, see <https://www.gnu.org/licenses/>.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAlbumsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "immich_albums" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.#", "10"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.album_name", "Test 53"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.album_thumbnail_asset_id", ""),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.description", "Test 53 Album"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.end_date", ""),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.has_shared_link", "false"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.is_activity_enabled", "true"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.order", "desc"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.shared", "false"),
					resource.TestCheckResourceAttr("data.immich_albums.test", "albums.0.start_date", ""),
				),
			},
		},
	})
}
