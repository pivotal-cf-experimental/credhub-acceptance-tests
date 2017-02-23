package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("updating a secret", func() {
	Describe("updating with a set (PUT)", func() {
		It("should be able to overwrite a secret", func() {
			credentialName := generateUniqueCredentialName()

			By("setting a new value secret", func() {
				session := runCommand("set", "-n", credentialName, "-t", "value", "-v", "old value")
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
				Expect(stdOut).To(MatchRegexp("Value:\\s+" + "old value"))
			})

			By("setting the value secret again", func() {
				session := runCommand("set", "-n", credentialName, "-t", "value", "-v", "new value")
				Eventually(session).Should(Exit(0))

				stdOut := string(session.Out.Contents())
				Expect(stdOut).To(MatchRegexp(`Type:\s+value`))
				Expect(stdOut).To(MatchRegexp("Value:\\s+" + "new value"))
			})
		})
	})

	Describe("generating -> setting -> generating", func() {
		It("does not bleed values from the generate", func() {
			caName := generateUniqueCredentialName()
			credentialname := generateUniqueCredentialName()

			By("generating a new ca", func() {
				runCommand("generate", "-n", caName, "-t", "certificate", "-c", "anything", "--is-ca", "--self-sign")
			})

			By("generating a new certificate signed by the CA", func() {
				session := runCommand("generate", "-n", credentialname, "-t", "certificate", "-c", "bla", "--ca", caName)
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`Ca:\s+-----BEGIN CERTIFICATE-----`))
				Expect(stdOut).To(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`))
			})

			By("overwriting the certificate with `set`", func() {
				session := runCommand("set", "-n", credentialname, "-t", "certificate", "--certificate-string", "fake-certificate")
				stdOut := string(session.Out.Contents())
				Eventually(session).Should(Exit(0))
				Expect(stdOut).To(MatchRegexp(`Type:\s+certificate`))
				Expect(stdOut).To(MatchRegexp(`Certificate:\s+fake-certificate`))
				Expect(stdOut).To(Not(MatchRegexp(`Ca:\s+-----BEGIN CERTIFICATE-----`)))
				Expect(stdOut).To(Not(MatchRegexp(`Private Key:\s+-----BEGIN RSA PRIVATE KEY-----`)))
			})
		})
	})
})