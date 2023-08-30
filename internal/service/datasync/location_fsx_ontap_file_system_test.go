// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datasync_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/datasync"
	"github.com/aws/aws-sdk-go/service/fsx"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfdatasync "github.com/hashicorp/terraform-provider-aws/internal/service/datasync"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccDataSyncLocationFSxONTAPFileSystem_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var v datasync.DescribeLocationFsxOntapOutput
	resourceName := "aws_datasync_location_fsx_ontap_file_system.test"
	fsResourceName := "aws_fsx_ontap_file_system.test"
	svmResourceName := "aws_fsx_ontap_storage_virtual_machine.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, fsx.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationFSxONTAPDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationFSxONTAPFileSystemConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationFSxONTAPExists(ctx, resourceName, &v),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "datasync", regexp.MustCompile(`location/loc-.+`)),
					resource.TestCheckResourceAttrSet(resourceName, "creation_time"),
					resource.TestCheckResourceAttrPair(resourceName, "fsx_filesystem_arn", fsResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_virtual_machine_arn", svmResourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestMatchResourceAttr(resourceName, "uri", regexp.MustCompile(`^fsxn-(nfs|smb)://.+/`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLocationFSxONTAPImportStateID(resourceName),
			},
		},
	})
}

func TestAccDataSyncLocationFSxONTAPFileSystem_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var v datasync.DescribeLocationFsxOntapOutput
	resourceName := "aws_datasync_location_fsx_ontap_file_system.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, fsx.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationFSxONTAPDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationFSxONTAPFileSystemConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationFSxONTAPExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfdatasync.ResourceLocationFSxONTAPFileSystem(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDataSyncLocationFSxONTAPFileSystem_subdirectory(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var v datasync.DescribeLocationFsxOntapOutput
	resourceName := "aws_datasync_location_fsx_ontap_file_system.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, fsx.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationFSxONTAPDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationFSxONTAPFileSystemConfig_subdirectory(rName, "/subdirectory1/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationFSxONTAPExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "subdirectory", "/subdirectory1/"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLocationFSxONTAPImportStateID(resourceName),
			},
		},
	})
}

func TestAccDataSyncLocationFSxONTAPFileSystem_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	var v datasync.DescribeLocationFsxOntapOutput
	resourceName := "aws_datasync_location_fsx_ontap_file_system.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, fsx.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, datasync.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckLocationFSxONTAPDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationFSxONTAPFileSystemConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationFSxONTAPExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLocationFSxONTAPImportStateID(resourceName),
			},
			{
				Config: testAccLocationFSxONTAPFileSystemConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationFSxONTAPExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccLocationFSxONTAPFileSystemConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLocationFSxONTAPExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func testAccCheckLocationFSxONTAPDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).DataSyncConn(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_datasync_location_fsx_ontap_file_system" {
				continue
			}

			_, err := tfdatasync.FindLocationFSxONTAPByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("DataSync Location FSx for NetApp ONTAP File System %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckLocationFSxONTAPExists(ctx context.Context, n string, v *datasync.DescribeLocationFsxOntapOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).DataSyncConn(ctx)

		output, err := tfdatasync.FindLocationFSxONTAPByARN(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccLocationFSxONTAPImportStateID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		return fmt.Sprintf("%s#%s", rs.Primary.ID, rs.Primary.Attributes["fsx_filesystem_arn"]), nil
	}
}

func testAccFSxOntapFileSystemConfig_base(rName string) string {
	return acctest.ConfigCompose(acctest.ConfigVPCWithSubnets(rName, 1), fmt.Sprintf(`
resource "aws_security_group" "test" {
  name   = %[1]q
  vpc_id = aws_vpc.test.id

  ingress {
    cidr_blocks = [aws_vpc.test.cidr_block]
    from_port   = 0
    protocol    = -1
    to_port     = 0
  }

  egress {
    cidr_blocks = ["0.0.0.0/0"]
    from_port   = 0
    protocol    = "-1"
    to_port     = 0
  }

  tags = {
    Name = %[1]q
  }
}

resource "aws_fsx_ontap_file_system" "test" {
  storage_capacity    = 1024
  subnet_ids          = aws_subnet.test[*].id
  deployment_type     = "SINGLE_AZ_1"
  throughput_capacity = 512
  preferred_subnet_id = aws_subnet.test[0].id
}

resource "aws_fsx_ontap_storage_virtual_machine" "test" {
  file_system_id = aws_fsx_ontap_file_system.test.id
  name           = %[1]q
}
`, rName))
}

func testAccLocationFSxONTAPFileSystemConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccFSxOntapFileSystemConfig_base(rName), `
resource "aws_datasync_location_fsx_ontap_file_system" "test" {
  security_group_arns         = [aws_security_group.test.arn]
  storage_virtual_machine_arn = aws_fsx_ontap_storage_virtual_machine.test.arn

  protocol {
    nfs {
      mount_options {
        version = "NFS3"
      }
    }
  }
}
`)
}

func testAccLocationFSxONTAPFileSystemConfig_subdirectory(rName, subdirectory string) string {
	return acctest.ConfigCompose(testAccFSxOntapFileSystemConfig_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_fsx_ontap_file_system" "test" {
  security_group_arns         = [aws_security_group.test.arn]
  storage_virtual_machine_arn = aws_fsx_ontap_storage_virtual_machine.test.arn
  subdirectory                = %[1]q

  protocol {
    nfs {
      mount_options {
        version = "NFS3"
      }
    }
  }
}
`, subdirectory))
}

func testAccLocationFSxONTAPFileSystemConfig_tags1(rName, key1, value1 string) string {
	return acctest.ConfigCompose(testAccFSxOntapFileSystemConfig_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_fsx_ontap_file_system" "test" {
  security_group_arns         = [aws_security_group.test.arn]
  storage_virtual_machine_arn = aws_fsx_ontap_storage_virtual_machine.test.arn

  tags = {
    %[1]q = %[2]q
  }

  protocol {
    nfs {
      mount_options {
        version = "NFS3"
      }
    }
  }
}
`, key1, value1))
}

func testAccLocationFSxONTAPFileSystemConfig_tags2(rName, key1, value1, key2, value2 string) string {
	return acctest.ConfigCompose(testAccFSxOntapFileSystemConfig_base(rName), fmt.Sprintf(`
resource "aws_datasync_location_fsx_ontap_file_system" "test" {
  security_group_arns         = [aws_security_group.test.arn]
  storage_virtual_machine_arn = aws_fsx_ontap_storage_virtual_machine.test.arn

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }

  protocol {
    nfs {
      mount_options {
        version = "NFS3"
      }
    }
  }
}
`, key1, value1, key2, value2))
}
