package integration_test

import (
	"github.com/cloudfoundry/libbuildpack/cutlass"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("running supply nodejs buildpack before the go buildpack", func() {
	var app *cutlass.App

	AfterEach(func() {
		if app != nil {
			app.Destroy()
		}
		app = nil
	})

	Context("the app is pushed", func() {
		BeforeEach(func() {
			if ok, err := cutlass.ApiGreaterThan("2.65.1"); err != nil || !ok {
				Skip("API version does not have multi-buildpack support")
			}

			app = cutlass.New(Fixtures("go_calls_nodejs"))
			app.Buildpacks = []string{
				"https://github.com/cloudfoundry/nodejs-buildpack#master",
				"go_buildpack",
			}
		})

		It("finds the supplied dependency in the runtime container", func() {
			PushAppAndConfirm(app)

			Expect(app.Stdout.String()).To(ContainSubstring("Nodejs Buildpack version"))

			Expect(app.GetBody("/")).To(MatchRegexp("INFO hello world"))
		})
	})
})
