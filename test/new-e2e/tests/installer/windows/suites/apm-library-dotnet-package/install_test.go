// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package dotnettests

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	winawshost "github.com/DataDog/datadog-agent/test/new-e2e/pkg/provisioners/aws/host/windows"
	installer "github.com/DataDog/datadog-agent/test/new-e2e/tests/installer/unix"
	installerwindows "github.com/DataDog/datadog-agent/test/new-e2e/tests/installer/windows"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/installer/windows/consts"
	"github.com/DataDog/datadog-agent/test/new-e2e/tests/windows"

	"testing"
)

var (
	//go:embed resources/web.config
	webConfigFile []byte
	//go:embed resources/index.aspx
	aspxFile []byte
)

type testDotnetLibraryInstallSuite struct {
	installerwindows.BaseSuite
}

// TestDotnetInstalls tests the usage of the Datadog installer to install the apm-library-dotnet-package package.
func TestDotnetLibraryInstalls(t *testing.T) {
	t.Parallel()
	e2e.Run(t, &testDotnetLibraryInstallSuite{},
		e2e.WithProvisioner(
			winawshost.ProvisionerNoAgentNoFakeIntake(
				winawshost.WithInstaller(),
			)))
}

// TestInstallUninstallDotnetLibraryPackage tests installing and uninstalling the Datadog APM Library for .NET using the Datadog installer.
func (s *testDotnetLibraryInstallSuite) TestInstallUninstallDotnetLibraryPackage() {
	s.Require().NoError(s.Installer().Install())
	defer s.Installer().Purge()
	s.installIIS()

	// TODO remove override once image is published
	output, err := s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to install the dotnet library package: %s", output)

	output, err = s.Installer().RemovePackage("datadog-apm-library-dotnet")

	s.Require().NoErrorf(err, "failed to remove the dotnet library package: %s", output)
	s.Require().Host(s.Env().RemoteHost).
		NoDirExists(consts.GetStableDirFor("datadog-apm-library-dotnet"),
			"the package directory should be removed")
}

// TestInstallUninstallDotnetLibraryPackage tests installing and uninstalling the Datadog APM Library for .NET using the Datadog installer.
func (s *testDotnetLibraryInstallSuite) TestReinstall() {
	s.Require().NoError(s.Installer().Install())
	defer s.Installer().Purge()
	s.installIIS()

	// TODO remove override once image is published
	output, err := s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to instal the dotnet library package: %s", output)

	output, err = s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to install the dotnet library package: %s", output)
}

func (s *testDotnetLibraryInstallSuite) TestUpdate() {
	s.Require().NoError(s.Installer().Install())
	defer s.Installer().Purge()
	s.installIIS()

	// TODO remove override once image is published
	output, err := s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to instal the dotnet library package: %s", output)

	output, err = s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to install the dotnet library package: %s", output)
}

func (s *testDotnetLibraryInstallSuite) TestRemovePackageFailsIfInUse() {
	s.Require().NoError(s.Installer().Install())
	defer s.Installer().Purge()
	s.installIIS()
	s.installAspNet()

	// TODO remove override once image is published
	output, err := s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to install the dotnet library package: %s", output)

	err = s.startIISApp()
	s.Require().NoError(err, "failed to start IIS app")

	output, err = s.Installer().RemovePackage("datadog-apm-library-dotnet")
	s.Require().Error(err, "Removing the package while the native profiler is used by another process should fail")

	err = s.stopIISApp()
	s.Require().NoError(err, "failed to stop IIS app")

	output, err = s.Installer().RemovePackage("datadog-apm-library-dotnet")
	s.Require().NoErrorf(err, "failed to remove the dotnet library package: %s", output)
}

func (s *testDotnetLibraryInstallSuite) TestCorruptedPackageGetsDeleted() {
	s.Require().NoError(s.Installer().Install())
	defer s.Installer().Purge()
	s.installIIS()

	// TODO remove override once image is published
	output, err := s.Installer().InstallPackage("datadog-apm-library-dotnet",
		installer.WithVersion("428c2fc49dc8e75040934d590fa52912f768ded7"),
		installer.WithRegistry("installtesting.datad0g.com"),
	)
	s.Require().NoErrorf(err, "failed to install the dotnet library package: %s", output)

	s.Env().RemoteHost.Remove(filepath.Join(consts.GetStableDirFor("datadog-apm-library-dotnet"), "installer", "Datadog.FleetInstaller.exe"))

	output, err = s.Installer().RemovePackage("datadog-apm-library-dotnet")
	s.Require().NoErrorf(err, "failed to remove the dotnet library package: %s", output)
}

func (s *testDotnetLibraryInstallSuite) installIIS() {
	s.BaseSuite.SetupSuite()
	host := s.Env().RemoteHost
	err := windows.InstallIIS(host)
	s.Require().NoError(err)
}

func (s *testDotnetLibraryInstallSuite) installAspNet() {
	host := s.Env().RemoteHost
	output, err := host.Execute("Install-WindowsFeature Web-Asp-Net45")
	s.Require().NoErrorf(err, "failed to install Asp.Net: %s", output)
}

func (s *testDotnetLibraryInstallSuite) startIISApp() error {
	host := s.Env().RemoteHost
	err := host.MkdirAll("C:\\inetpub\\wwwroot\\DummyApp")
	if err != nil {
		return fmt.Errorf("failed to create site directory: %w", err)
	}
	_, err = host.WriteFile("C:\\inetpub\\wwwroot\\DummyApp\\web.config", webConfigFile)
	if err != nil {
		return fmt.Errorf("failed to write web.config file: %w", err)
	}
	_, err = host.WriteFile("C:\\inetpub\\wwwroot\\DummyApp\\index.aspx", aspxFile)
	if err != nil {
		return fmt.Errorf("failed to write web.config file: %w", err)
	}
	script := `
$SitePath = "C:\inetpub\wwwroot\DummyApp"
New-WebSite -Name DummyApp -PhysicalPath $SitePath -Port 8080 -ApplicationPool "DefaultAppPool" -Force
Stop-WebSite -Name "DummyApp"
Start-WebSite -Name "DummyApp"
$state = (Get-WebAppPoolState -Name "DefaultAppPool").Value
if ($state -eq "Stopped") {
    Start-WebAppPool -Name "DefaultAppPool"
}
Restart-WebAppPool -Name "DefaultAppPool"
Invoke-WebRequest -Uri "http://localhost:8080/index.aspx" -UseBasicParsing
	`
	output, err := host.Execute(script)
	if err != nil {
		return fmt.Errorf("failed to start site: %w\n%s", err, output)
	}
	return nil
}

func (s *testDotnetLibraryInstallSuite) stopIISApp() error {
	script := `
Stop-WebSite -Name "DummyApp"
Stop-WebAppPool -Name "DefaultAppPool"
$retryCount = 0
do {
    Start-Sleep -Seconds 1
    $status = (Get-WebAppPoolState -Name DefaultAppPool).Value
    $retryCount++
} while ($status -ne "Stopped" -and $retryCount -lt 60)

if ($status -ne "Stopped") {
	exit -1
}
	`
	host := s.Env().RemoteHost
	output, err := host.Execute(script)
	if err != nil {
		return fmt.Errorf("failed to start site: %w\n%s", err, output)
	}
	return nil
}
