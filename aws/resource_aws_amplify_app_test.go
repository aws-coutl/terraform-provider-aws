package aws

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/amplify"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/naming"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/amplify/finder"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/tfresource"
)

// TODO sweeper

func TestAccAWSAmplifyApp_basic(t *testing.T) {
	var app amplify.App
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigName(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					resource.TestCheckNoResourceAttr(resourceName, "access_token"),
					testAccMatchResourceAttrRegionalARN(resourceName, "arn", "amplify", regexp.MustCompile(`apps/.+`)),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_patterns.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "basic_auth_credentials", ""),
					resource.TestCheckResourceAttr(resourceName, "build_spec", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_headers", ""),
					resource.TestCheckResourceAttr(resourceName, "custom_rule.#", "0"),
					resource.TestMatchResourceAttr(resourceName, "default_domain", regexp.MustCompile(`\.amplifyapp\.com$`)),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_auto_branch_creation", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_basic_auth", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_branch_auto_build", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_branch_auto_deletion", "false"),
					resource.TestCheckResourceAttr(resourceName, "environment_variables.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "iam_service_role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", ""),
					resource.TestCheckNoResourceAttr(resourceName, "oauth_token"),
					resource.TestCheckResourceAttr(resourceName, "platform", "WEB"),
					resource.TestCheckResourceAttr(resourceName, "production_branch.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "repository", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAWSAmplifyApp_Name_Generated(t *testing.T) {
	var app amplify.App
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigNameGenerated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					naming.TestCheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "terraform-"),
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

func TestAccAWSAmplifyApp_NamePrefix(t *testing.T) {
	var app amplify.App
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigNamePrefix("tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					naming.TestCheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
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

func TestAccAWSAmplifyApp_Tags(t *testing.T) {
	var app amplify.App
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
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
				Config: testAccAWSAmplifyAppConfigTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSAmplifyAppConfigTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_rename(t *testing.T) {
	resourceName := "aws_amplify_app.test"

	// name is not unique and can be renamed
	rName1 := acctest.RandomWithPrefix("tf-acc-test")
	rName2 := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigName(rName1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigName(rName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_description(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	// once set, description cannot be removed.
	description1 := acctest.RandomWithPrefix("tf-acc-test")
	description2 := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigDescription(rName, description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigDescription(rName, description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", description2),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_repository(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigRepository(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "repository", regexp.MustCompile("^https://github.com")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// access_token is ignored because AWS does not store access_token and oauth_token
				// See https://docs.aws.amazon.com/sdk-for-go/api/service/amplify/#CreateAppInput
				ImportStateVerifyIgnore: []string{"access_token"},
			},
		},
	})
}

func TestAccAWSAmplifyApp_buildSpec(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	// once set, build_spec cannot be removed.
	buildSpec1 := "version: 0.1"
	buildSpec2 := "version: 0.2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigBuildSpec(rName, buildSpec1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "build_spec", buildSpec1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigBuildSpec(rName, buildSpec2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "build_spec", buildSpec2),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_customRules(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigCustomRules1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "custom_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.source", "/<*>"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.status", "404"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.target", "/index.html"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigCustomRules2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "custom_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.source", "/documents"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.status", "302"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.target", "/documents/us"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.0.condition", "<US>"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.1.source", "/<*>"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.1.status", "200"),
					resource.TestCheckResourceAttr(resourceName, "custom_rules.1.target", "/index.html"),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_environmentVariables(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigEnvironmentVariables1(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environment_variables.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "environment_variables.ENVVAR1", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigEnvironmentVariables2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environment_variables.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "environment_variables.ENVVAR1", "2"),
					resource.TestCheckResourceAttr(resourceName, "environment_variables.ENVVAR2", "2"),
				),
			},
			{
				Config: testAccAWSAmplifyAppConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environment_variables.%", "0"),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_autoBranchCreationConfig(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigAutoBranchCreationConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_auto_branch_creation", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.auto_branch_creation_patterns.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.auto_branch_creation_patterns.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.auto_branch_creation_patterns.1", "*/**"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.build_spec", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.framework", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.stage", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.basic_auth_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_auto_build", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_pull_request_preview", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.pull_request_environment_name", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.environment_variables.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigAutoBranchCreationConfigModified(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_auto_branch_creation", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.auto_branch_creation_patterns.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.auto_branch_creation_patterns.0", "feature/*"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.build_spec", "version: 0.1"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.framework", "React"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.stage", "DEVELOPMENT"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.basic_auth_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.basic_auth_config.0.enable_basic_auth", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.basic_auth_config.0.username", "username"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.basic_auth_config.0.password", "password"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_auto_build", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_pull_request_preview", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.pull_request_environment_name", "env"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.environment_variables.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.environment_variables.ENVVAR1", "1"),
				),
			},
			{
				Config: testAccAWSAmplifyAppConfigAutoBranchCreationConfigModified2(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.enable_auto_branch_creation", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.basic_auth_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.0.environment_variables.%", "0"),
				),
			},
			{
				Config: testAccAWSAmplifyAppConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "auto_branch_creation_config.#", "0"),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_BasicAuthCredentials(t *testing.T) {
	var app amplify.App
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	credentials1 := base64.StdEncoding.EncodeToString([]byte("username1:password1"))
	credentials2 := base64.StdEncoding.EncodeToString([]byte("username2:password2"))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigBasicAuthCredentials(rName, credentials1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "basic_auth_credentials", credentials1),
					resource.TestCheckResourceAttr(resourceName, "enable_basic_auth", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigBasicAuthCredentials(rName, credentials2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "basic_auth_credentials", credentials2),
					resource.TestCheckResourceAttr(resourceName, "enable_basic_auth", "true"),
				),
			},
			{
				Config: testAccAWSAmplifyAppConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAmplifyAppExists(resourceName, &app),
					resource.TestCheckResourceAttr(resourceName, "basic_auth_credentials", ""),
					resource.TestCheckResourceAttr(resourceName, "enable_basic_auth", "false"),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_enableBranchAutoBuild(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigEnableBranchAutoBuild(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_branch_auto_build", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigName(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enable_branch_auto_build", "false"),
				),
			},
		},
	})
}

func TestAccAWSAmplifyApp_iamServiceRoleArn(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_amplify_app.test"

	roleName1 := acctest.RandomWithPrefix("tf-acc-test")
	roleName2 := acctest.RandomWithPrefix("tf-acc-test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckAWSAmplify(t) },
		ErrorCheck:   testAccErrorCheck(t, amplify.EndpointsID),
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAWSAmplifyAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAmplifyAppConfigIAMServiceRoleArn(rName, roleName1),
				Check: resource.ComposeTestCheckFunc(
					testAccMatchResourceAttrGlobalARN(resourceName, "iam_service_role_arn", "iam", regexp.MustCompile("role/"+roleName1)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAmplifyAppConfigIAMServiceRoleArn(rName, roleName2),
				Check: resource.ComposeTestCheckFunc(
					testAccMatchResourceAttrGlobalARN(resourceName, "iam_service_role_arn", "iam", regexp.MustCompile("role/"+roleName2)),
				),
			},
		},
	})
}

func testAccCheckAWSAmplifyAppExists(n string, v *amplify.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Amplify App ID is set")
		}

		conn := testAccProvider.Meta().(*AWSClient).amplifyconn

		output, err := finder.AppByID(conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccCheckAWSAmplifyAppDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).amplifyconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_amplify_app" {
			continue
		}

		_, err := finder.AppByID(conn, rs.Primary.ID)

		if tfresource.NotFound(err) {
			continue
		}

		if err != nil {
			return err
		}

		return fmt.Errorf("Amplify App %s still exists", rs.Primary.ID)
	}

	return nil
}

func testAccPreCheckAWSAmplify(t *testing.T) {
	if testAccGetPartition() == "aws-us-gov" {
		t.Skip("AWS Amplify is not supported in GovCloud partition")
	}
}

func testAccAWSAmplifyAppConfigName(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = %[1]q
}
`, rName)
}

func testAccAWSAmplifyAppConfigNameGenerated() string {
	return `
resource "aws_amplify_app" "test" {}
`
}

func testAccAWSAmplifyAppConfigNamePrefix(namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}

func testAccAWSAmplifyAppConfigDescription(rName string, description string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  description = "%s"
}
`, rName, description)
}

func testAccAWSAmplifyAppConfigRepository(rName string) string {
	repository := os.Getenv("AMPLIFY_GITHUB_REPOSITORY")
	accessToken := os.Getenv("AMPLIFY_GITHUB_ACCESS_TOKEN")

	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  repository   = "%s"
  access_token = "%s"
}
`, rName, repository, accessToken)
}

func testAccAWSAmplifyAppConfigBuildSpec(rName string, buildSpec string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  build_spec = "%s"
}
`, rName, buildSpec)
}

func testAccAWSAmplifyAppConfigCustomRules1(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  custom_rules {
    source = "/<*>"
    status = "404"
    target = "/index.html"
  }
}
`, rName)
}

func testAccAWSAmplifyAppConfigCustomRules2(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  custom_rules {
    source    = "/documents"
    status    = "302"
    target    = "/documents/us"
    condition = "<US>"
  }

  custom_rules {
    source = "/<*>"
    status = "200"
    target = "/index.html"
  }
}
`, rName)
}

func testAccAWSAmplifyAppConfigEnvironmentVariables1(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  environment_variables = {
    ENVVAR1 = "1"
  }
}
`, rName)
}

func testAccAWSAmplifyAppConfigEnvironmentVariables2(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  environment_variables = {
    ENVVAR1 = "2",
    ENVVAR2 = "2"
  }
}
`, rName)
}

func testAccAWSAmplifyAppConfigAutoBranchCreationConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  auto_branch_creation_config {
    enable_auto_branch_creation = true

    auto_branch_creation_patterns = [
      "*",
      "*/**",
    ]
  }
}
`, rName)
}

func testAccAWSAmplifyAppConfigAutoBranchCreationConfigModified(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  auto_branch_creation_config {
    enable_auto_branch_creation = true

    auto_branch_creation_patterns = [
      "feature/*",
    ]

    build_spec = "version: 0.1"
    framework  = "React"
    stage      = "DEVELOPMENT"

    basic_auth_config {
      enable_basic_auth = true
      username          = "username"
      password          = "password"
    }

    enable_auto_build = true

    enable_pull_request_preview   = true
    pull_request_environment_name = "env"

    environment_variables = {
      ENVVAR1 = "1"
    }
  }
}

`, rName)
}

func testAccAWSAmplifyAppConfigAutoBranchCreationConfigModified2(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  auto_branch_creation_config {
    enable_auto_branch_creation = true

    auto_branch_creation_patterns = [
      "feature/*",
    ]

    build_spec = "version: 0.1"
    framework  = "React"
    stage      = "DEVELOPMENT"

    enable_auto_build = false

    enable_pull_request_preview   = false
    pull_request_environment_name = "env"
  }
}
`, rName)
}

func testAccAWSAmplifyAppConfigBasicAuthCredentials(rName, basicAuthCredentials string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = %[1]q

  basic_auth_credentials = %[2]q
  enable_basic_auth      = true
}
`, rName, basicAuthCredentials)
}

func testAccAWSAmplifyAppConfigEnableBranchAutoBuild(rName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  enable_branch_auto_build = true
}
`, rName)
}

func testAccAWSAmplifyAppConfigIAMServiceRoleArn(rName string, roleName string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = "%s"

  iam_service_role_arn = aws_iam_role.role.arn
}

resource "aws_iam_role" "role" {
  name = "%s"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "amplify.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}
`, rName, roleName)
}

func testAccAWSAmplifyAppConfigTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccAWSAmplifyAppConfigTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_amplify_app" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
