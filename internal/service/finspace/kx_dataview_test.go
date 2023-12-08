package finspace_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/finspace"
	"github.com/aws/aws-sdk-go-v2/service/finspace/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tffinspace "github.com/hashicorp/terraform-provider-aws/internal/service/finspace"
	"github.com/hashicorp/terraform-provider-aws/names"
	"testing"
)

func TestAccFinSpaceKxDataview_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	ctx := acctest.Context(t)
	var kxdataview finspace.GetKxDataviewOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_finspace_kx_dataview.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, finspace.ServiceID)
		},
		ErrorCheck:               acctest.ErrorCheck(t, finspace.ServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckKxDataviewDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccKxDataviewConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKxDataviewExists(ctx, resourceName, &kxdataview),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", string(types.KxDataviewStatusActive)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccKxDataviewConfigBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_kms_key" "test" {
  deletion_window_in_days = 7
}

resource "aws_finspace_kx_environment" "test" {
	name       = %[1]q
    kms_key_id = aws_kms_key.test.arn
}
resource "aws_finspace_kx_database" "test" {
  name           = %[1]q
  environment_id = aws_finspace_kx_environment.test.id
}
`, rName)
}
func testAccKxDataviewConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccKxDataviewConfigBase(rName),
		fmt.Sprintf(`
resource "aws_finspace_kx_dataview" "test" {
  name                 = %[1]q
  environment_id       = aws_finspace_kx_environment.test.id
  database_name        = aws_finspace_kx_database.test.name
  auto_update          = true
  az_mode              = "SINGLE"
  availability_zone_id = aws_finspace_kx_environment.test.availability_zones[0]
}
`, rName))
}

func testAccCheckKxDataviewExists(ctx context.Context, name string, kxdataview *finspace.GetKxDataviewOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.FinSpace, create.ErrActionCheckingExistence, tffinspace.ResNameKxDataview, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.FinSpace, create.ErrActionCheckingExistence, tffinspace.ResNameKxDataview, name, errors.New("not set"))
		}

		conn := tffinspace.TempFinspaceClient()
		resp, err := conn.GetKxDataview(ctx, &finspace.GetKxDataviewInput{
			DatabaseName:  aws.String(rs.Primary.Attributes["database_name"]),
			EnvironmentId: aws.String(rs.Primary.Attributes["environment_id"]),
			DataviewName:  aws.String(rs.Primary.Attributes["name"]),
		})
		if err != nil {
			return create.Error(names.FinSpace, create.ErrActionCheckingExistence, tffinspace.ResNameKxDataview, rs.Primary.ID, err)
		}

		*kxdataview = *resp

		return nil
	}
}

func testAccCheckKxDataviewDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_finspace_kx_dataview" {
				continue
			}

			conn := tffinspace.TempFinspaceClient()
			_, err := conn.GetKxDataview(ctx, &finspace.GetKxDataviewInput{
				DatabaseName:  aws.String(rs.Primary.Attributes["database_name"]),
				EnvironmentId: aws.String(rs.Primary.Attributes["environment_id"]),
				DataviewName:  aws.String(rs.Primary.Attributes["name"]),
			})
			if err != nil {
				var nfe *types.ResourceNotFoundException
				if errors.As(err, &nfe) {
					return nil
				}
				return err
			}
			return create.Error(names.FinSpace, create.ErrActionCheckingExistence, tffinspace.ResNameKxDataview, rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccKxDataviewVolumeBase(rName string) string {
	return fmt.Sprintf(`
resource "aws_finspace_kx_volume" "test" {
  name                 = %[1]q
  environment_id       = aws_finspace_kx_environment.test.id
  availability_zones = [aws_finspace_kx_environment.test.availability_zones[0]]
  az_mode              = "SINGLE"
  type                 = "NAS_1"
  nas1_configuration {
	  type= "SSD_250"
	  size= 1200
  }
}
`, rName)
}

func testAccKxDataviewConfig_withKxVolume(rName string) string {
	return acctest.ConfigCompose(
		testAccKxDataviewConfigBase(rName),
		testAccKxDataviewVolumeBase(rName),
		fmt.Sprintf(`
resource "aws_finspace_kx_dataview" "test" {
  name                 = %[1]q
  environment_id       = aws_finspace_kx_environment.test.id
  database_name        = aws_finspace_kx_database.test.name
  auto_update          = true
  az_mode              = "SINGLE"
  availability_zone_id = aws_finspace_kx_environment.test.availability_zones[0]

  segment_configurations {
      db_paths = ["/*"]
 	  volume_name = aws_finspace_kx_volume.test.name
  }
}
`, rName))
}

func TestAccFinSpaceKxDataview_withKxVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}
	ctx := acctest.Context(t)

	var kxdataview finspace.GetKxDataviewOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_finspace_kx_dataview.test"

	resource.ParallelTest(t, resource.TestCase{

		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, finspace.ServiceID)
		},
		ErrorCheck: acctest.ErrorCheck(t, finspace.ServiceID),

		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,

		CheckDestroy: testAccCheckKxDataviewDestroy(ctx),

		Steps: []resource.TestStep{
			{
				Config: testAccKxDataviewConfig_withKxVolume(rName),

				Check: resource.ComposeTestCheckFunc(
					testAccCheckKxDataviewExists(ctx, resourceName, &kxdataview),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", string(types.KxDataviewStatusActive)),
				),
			},
		},
	})
}
