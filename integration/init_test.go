package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

var settings struct {
	Buildpacks struct {
		BundleInstall struct {
			Online string
		}
		Bundler struct {
			Online string
		}
		MRI struct {
			Online string
		}
		NodeEngine struct {
			Online string
		}
		Puma struct {
			Online string
		}
		RailsAssets struct {
			Online string
		}
		Yarn struct {
			Online string
		}
		YarnInstall struct {
			Online string
		}
		NpmInstall struct {
			Online string
		}
	}

	Buildpack struct {
		ID   string
		Name string
	}

	Config struct {
		BundleInstall string `json:"bundle-install"`
		Bundler       string `json:"bundler"`
		MRI           string `json:"mri"`
		NodeEngine    string `json:"node-engine"`
		Puma          string `json:"puma"`
		Yarn          string `json:"yarn"`
		YarnInstall   string `json:"yarn-install"`
		NpmInstall    string `json:"npm-install"`
	}

	Pack   occam.Pack
	Docker occam.Docker
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect
	SetDefaultEventuallyTimeout(30 * time.Second)
	format.MaxLength = 0

	root, err := filepath.Abs("./..")
	Expect(err).NotTo(HaveOccurred())

	file, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())

	file, err = os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&settings)
	Expect(err).NotTo(HaveOccurred())
	Expect(file.Close()).To(Succeed())

	buildpackStore := occam.NewBuildpackStore()

	settings.Buildpacks.RailsAssets.Online, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.BundleInstall.Online, err = buildpackStore.Get.
		Execute(settings.Config.BundleInstall)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Bundler.Online, err = buildpackStore.Get.
		Execute(settings.Config.Bundler)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.MRI.Online, err = buildpackStore.Get.
		Execute(settings.Config.MRI)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.NodeEngine.Online, err = buildpackStore.Get.
		Execute(settings.Config.NodeEngine)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Puma.Online, err = buildpackStore.Get.
		Execute(settings.Config.Puma)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Yarn.Online, err = buildpackStore.Get.
		Execute(settings.Config.Yarn)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.YarnInstall.Online, err = buildpackStore.Get.
		Execute(settings.Config.YarnInstall)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.NpmInstall.Online, err = buildpackStore.Get.
		Execute(settings.Config.NpmInstall)
	Expect(err).NotTo(HaveOccurred())

	settings.Pack = occam.NewPack().WithVerbose()
	settings.Docker = occam.NewDocker()

	suite := spec.New("Integration", spec.Parallel(), spec.Report(report.Terminal{}))
	suite("Rails6.1", testRails61)
	suite("Rails7.0", testRails70)
	suite("ReusingLayerRebuild", testReusingLayerRebuild)
	suite.Run(t)
}
