package apprunner_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apprunner"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfapprunner "github.com/hashicorp/terraform-provider-aws/internal/service/apprunner"
)

func TestAccAppRunnerAutoScalingConfigurationVersion_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_auto_scaling_configuration_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutoScalingConfigurationVersionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingConfigurationVersionConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_concurrency", "100"),
					resource.TestCheckResourceAttr(resourceName, "max_size", "25"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
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

func TestAccAppRunnerAutoScalingConfigurationVersion_complex(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_auto_scaling_configuration_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutoScalingConfigurationVersionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingConfigurationVersionConfig_nonDefaults(rName, 50, 10, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_concurrency", "50"),
					resource.TestCheckResourceAttr(resourceName, "max_size", "10"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test resource recreation such that the revision number is still 1
				Config: testAccAutoScalingConfigurationVersionConfig_nonDefaults(rName, 150, 20, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_concurrency", "150"),
					resource.TestCheckResourceAttr(resourceName, "max_size", "20"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test resource recreation such that the revision number is still 1
				Config: testAccAutoScalingConfigurationVersionConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_concurrency", "100"),
					resource.TestCheckResourceAttr(resourceName, "max_size", "25"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
				),
			},
		},
	})
}

func TestAccAppRunnerAutoScalingConfigurationVersion_multipleVersions(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_auto_scaling_configuration_version.test"
	otherResourceName := "aws_apprunner_auto_scaling_configuration_version.other"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutoScalingConfigurationVersionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingConfigurationVersionConfig_multipleVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					testAccCheckAutoScalingConfigurationVersionExists(ctx, otherResourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest", "true"),
					resource.TestCheckResourceAttr(resourceName, "max_concurrency", "100"),
					resource.TestCheckResourceAttr(resourceName, "max_size", "25"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
					acctest.MatchResourceAttrRegionalARN(otherResourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/2/.+`, rName))),
					resource.TestCheckResourceAttr(otherResourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(otherResourceName, "auto_scaling_configuration_revision", "2"),
					resource.TestCheckResourceAttr(otherResourceName, "latest", "true"),
					resource.TestCheckResourceAttr(otherResourceName, "max_concurrency", "100"),
					resource.TestCheckResourceAttr(otherResourceName, "max_size", "25"),
					resource.TestCheckResourceAttr(otherResourceName, "min_size", "1"),
					resource.TestCheckResourceAttr(otherResourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
				),
			},
			{
				// Test update of "latest" computed attribute after apply
				Config: testAccAutoScalingConfigurationVersionConfig_multipleVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					testAccCheckAutoScalingConfigurationVersionExists(ctx, otherResourceName),
					resource.TestCheckResourceAttr(resourceName, "latest", "false"),
					resource.TestCheckResourceAttr(otherResourceName, "latest", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      otherResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAppRunnerAutoScalingConfigurationVersion_updateMultipleVersions(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_auto_scaling_configuration_version.test"
	otherResourceName := "aws_apprunner_auto_scaling_configuration_version.other"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutoScalingConfigurationVersionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingConfigurationVersionConfig_multipleVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					testAccCheckAutoScalingConfigurationVersionExists(ctx, otherResourceName),
				),
			},
			{
				Config: testAccAutoScalingConfigurationVersionConfig_updateMultipleVersions(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					testAccCheckAutoScalingConfigurationVersionExists(ctx, otherResourceName),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/1/.+`, rName))),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_configuration_revision", "1"),
					resource.TestCheckResourceAttr(resourceName, "latest", "false"),
					resource.TestCheckResourceAttr(resourceName, "max_concurrency", "100"),
					resource.TestCheckResourceAttr(resourceName, "max_size", "25"),
					resource.TestCheckResourceAttr(resourceName, "min_size", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
					acctest.MatchResourceAttrRegionalARN(otherResourceName, "arn", "apprunner", regexp.MustCompile(fmt.Sprintf(`autoscalingconfiguration/%s/2/.+`, rName))),
					resource.TestCheckResourceAttr(otherResourceName, "auto_scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(otherResourceName, "auto_scaling_configuration_revision", "2"),
					resource.TestCheckResourceAttr(otherResourceName, "latest", "true"),
					resource.TestCheckResourceAttr(otherResourceName, "max_concurrency", "125"),
					resource.TestCheckResourceAttr(otherResourceName, "max_size", "20"),
					resource.TestCheckResourceAttr(otherResourceName, "min_size", "1"),
					resource.TestCheckResourceAttr(otherResourceName, "status", tfapprunner.AutoScalingConfigurationStatusActive),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      otherResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAppRunnerAutoScalingConfigurationVersion_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_auto_scaling_configuration_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutoScalingConfigurationVersionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingConfigurationVersionConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfapprunner.ResourceAutoScalingConfigurationVersion(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAppRunnerAutoScalingConfigurationVersion_tags(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_apprunner_auto_scaling_configuration_version.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, apprunner.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutoScalingConfigurationVersionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoScalingConfigurationVersionConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutoScalingConfigurationVersionConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAutoScalingConfigurationVersionConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutoScalingConfigurationVersionExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAutoScalingConfigurationVersionDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_apprunner_auto_scaling_configuration_version" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).AppRunnerConn()

			input := &apprunner.DescribeAutoScalingConfigurationInput{
				AutoScalingConfigurationArn: aws.String(rs.Primary.ID),
			}

			output, err := conn.DescribeAutoScalingConfigurationWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, apprunner.ErrCodeResourceNotFoundException) {
				continue
			}

			if err != nil {
				return err
			}

			if output != nil && output.AutoScalingConfiguration != nil && aws.StringValue(output.AutoScalingConfiguration.Status) != "inactive" {
				return fmt.Errorf("App Runner AutoScaling Configuration (%s) still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckAutoScalingConfigurationVersionExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Runner Service ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppRunnerConn()

		input := &apprunner.DescribeAutoScalingConfigurationInput{
			AutoScalingConfigurationArn: aws.String(rs.Primary.ID),
		}

		output, err := conn.DescribeAutoScalingConfigurationWithContext(ctx, input)

		if err != nil {
			return err
		}

		if output == nil || output.AutoScalingConfiguration == nil {
			return fmt.Errorf("App Runner AutoScaling Configuration (%s) not found", rs.Primary.ID)
		}

		return nil
	}
}

func testAccAutoScalingConfigurationVersionConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_apprunner_auto_scaling_configuration_version" "test" {
  auto_scaling_configuration_name = %[1]q
}
`, rName)
}

func testAccAutoScalingConfigurationVersionConfig_nonDefaults(rName string, maxConcurrency, maxSize, minSize int) string {
	return fmt.Sprintf(`
resource "aws_apprunner_auto_scaling_configuration_version" "test" {
  auto_scaling_configuration_name = %[1]q

  max_concurrency = %[2]d
  max_size        = %[3]d
  min_size        = %[4]d
}
`, rName, maxConcurrency, maxSize, minSize)
}

func testAccAutoScalingConfigurationVersionConfig_multipleVersions(rName string) string {
	return fmt.Sprintf(`
resource "aws_apprunner_auto_scaling_configuration_version" "test" {
  auto_scaling_configuration_name = %[1]q
}

resource "aws_apprunner_auto_scaling_configuration_version" "other" {
  auto_scaling_configuration_name = aws_apprunner_auto_scaling_configuration_version.test.auto_scaling_configuration_name
}
`, rName)
}

func testAccAutoScalingConfigurationVersionConfig_updateMultipleVersions(rName string) string {
	return fmt.Sprintf(`
resource "aws_apprunner_auto_scaling_configuration_version" "test" {
  auto_scaling_configuration_name = %[1]q
}

resource "aws_apprunner_auto_scaling_configuration_version" "other" {
  auto_scaling_configuration_name = aws_apprunner_auto_scaling_configuration_version.test.auto_scaling_configuration_name

  max_concurrency = 125
  max_size        = 20
}
`, rName)
}

func testAccAutoScalingConfigurationVersionConfig_tags1(rName string, tagKey1 string, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_apprunner_auto_scaling_configuration_version" "test" {
  auto_scaling_configuration_name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccAutoScalingConfigurationVersionConfig_tags2(rName string, tagKey1 string, tagValue1 string, tagKey2 string, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_apprunner_auto_scaling_configuration_version" "test" {
  auto_scaling_configuration_name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
