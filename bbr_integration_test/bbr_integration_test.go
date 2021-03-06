package bbr_integration

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-acceptance-tests/test_helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Backup and Restore", func() {
	var credentialName string
	var bbrTestPath = "bbr_test"

	BeforeEach(func() {
		credentialName = fmt.Sprintf("%s/%s", bbrTestPath, test_helpers.GenerateUniqueCredentialName())

		By("authenticating against credhub")
		RunCommand("credhub", "api", "--server", config.ApiUrl, "--skip-tls-validation")
		RunCommand("credhub", "login", "--skip-tls-validation", "-u", config.ApiUsername, "-p", config.ApiPassword)

		CleanupCredhub(bbrTestPath)
	})

	AfterEach(func() {
		CleanupCredhub(bbrTestPath)
	})

	It("Successfully backs up and restores a Credhub release", func() {
		By("adding a test credential")
		RunCommand("credhub", "set", "--name", credentialName, "--value", "originalsecret")

		By("running bbr backup")
		RunCommand("bbr", "deployment", "--target", config.Bosh.URL, "--ca-cert", config.Bosh.CertPath, "--username",
			config.Bosh.Client, "--password", config.Bosh.ClientSecret, "--deployment", config.Bosh.DeploymentName, "backup")

		By("asserting that the backup archive exists and contains a pg dump file")
		RunCommand("tar", "zxvf", config.Bosh.DeploymentName+"/credhub-0.tgz")
		Eventually(RunCommand("ls", "./credhub/credhubdb_dump")).Should(gexec.Exit(0))

		By("editing the test credential")
		RunCommand("credhub", "set", "--name", credentialName, "--value", "updatedsecret")

		By("running bbr restore")
		RunCommand("bbr", "deployment", "--target", config.Bosh.URL, "--ca-cert", config.Bosh.CertPath, "--username",
			config.Bosh.Client, "--password", config.Bosh.ClientSecret, "--deployment", config.Bosh.DeploymentName, "restore")

		By("checking if the test credentials was restored")
		getSession := RunCommand("credhub", "get", "--name", credentialName)
		Eventually(getSession).Should(gexec.Exit(0))
		Eventually(getSession.Out).Should(gbytes.Say("value: originalsecret"))
	})
})

func CleanupCredhub(path string) {
	By("Cleaning up credhub bbr test passwords")
	RunCommand(
		"sh", "-c",
		fmt.Sprintf("credhub find -p /%s | tail -n +2 | cut -d\" \" -f1 | xargs -IN credhub delete --name N", path),
	)
}
