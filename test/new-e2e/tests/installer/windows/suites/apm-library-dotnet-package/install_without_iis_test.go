package dotnettests

import (
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	winawshost "github.com/DataDog/datadog-agent/test/new-e2e/pkg/provisioners/aws/host/windows"
	installer "github.com/DataDog/datadog-agent/test/new-e2e/tests/installer/unix"
	installerwindows "github.com/DataDog/datadog-agent/test/new-e2e/tests/installer/windows"

	"testing"
)

type testDotnetLibraryInstallSuiteWithoutIIS struct {
	installerwindows.BaseSuite
}

// TestDotnetInstalls tests the usage of the Datadog installer to install the apm-library-dotnet-package package.
func TestDotnetLibraryInstallsWithoutIIS(t *testing.T) {
	e2e.Run(t, &testDotnetLibraryInstallSuiteWithoutIIS{},
		e2e.WithProvisioner(winawshost.ProvisionerNoAgentNoFakeIntake(
			winawshost.WithInstaller(),
		)))
}

// TestInstallDotnetLibraryPackageWithoutIIS tests installing the Datadog APM Library for .NET using the Datadog installer without IIS installed.
func (s *testDotnetLibraryInstallSuiteWithoutIIS) TestInstallDotnetLibraryPackageWithoutIIS() {
	s.Require().NoError(s.Installer().Install())
	defer s.Installer().Purge()

	// TODO:DONOTMERGE remove override once image is published
	_, err := s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().Error(err, "Installing the dotnet library package without IIS should fail")
	// TODO check that the package gets deleted
	// s.Require().Host(s.Env().RemoteHost).
	// 	NoDirExists(consts.GetStableDirFor("datadog-apm-library-dotnet"),
	// 		"the package directory should not exist")
}
