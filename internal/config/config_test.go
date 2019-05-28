package config

import (
	"testing"
)

func TestPrivateloadConfig(t *testing.T) {
	ec := Configuration{"goProbe", "localhost", 5005, "C:\\Logs\\error.log", false, 120}
	ac := loadConfig()

	if ac.GetApplicationName() != ec.GetApplicationName() {
		t.Errorf("Invalid Configuration Application Name, got : %s want : %s", ac.GetApplicationName(), ec.GetApplicationName())
	}
	if ac.GetHostName() != ec.GetHostName() {
		t.Errorf("Invalid Configuration Host Name, got : %s want : %s", ac.GetHostName(), ec.GetHostName())
	}
	if ac.GetHostPort() != ec.GetHostPort() {
		t.Errorf("Invalid Configuration Host Port, got : %d want : %d", ac.GetHostPort(), ec.GetHostPort())
	}
	if ac.GetLogPath() != ec.GetLogPath() {
		t.Errorf("Invalid Configuration Log Path, got : %s want : %s", ac.GetLogPath(), ec.GetLogPath())
	}
	if ac.GetDebugEnabled() != ec.GetDebugEnabled() {
		t.Errorf("Invalid Configuration Debug Enabled, got : %t want : %t", ac.GetDebugEnabled(), ec.GetDebugEnabled())
	}
	if ac.GetInterval() != ec.GetInterval() {
		t.Errorf("Invalid Configuration Interval, got : %d want : %d", ac.GetInterval(), ec.GetInterval())
	}
}

func TestPrivatevalidate(t *testing.T) {

	ec := Configuration{
		HostName: "localhost",
		HostPort: 5005,
	}
	loadLogger(&ec)
	validate(&ec)

	if ec.GetHostName() != "localhost" {
		t.Errorf("Invalid Host Name, got : %s want : %s", ec.GetHostName(), "localhost")
	}
	if ec.GetHostPort() != 5005 {
		t.Errorf("Invalid Host Port, got : %d want : %d", ec.GetHostPort(), 5005)
	}
	if ec.GetApplicationName() != "goProbe" {
		t.Errorf("Invalid Application Name, got : %s want : %s", ec.GetApplicationName(), "goProbe")
	}
	if ec.GetInterval() != 120 {
		t.Errorf("Invalid Interval, got : %d want : %d", ec.GetHostPort(), 120)
	}
}

func TestPrivateloadLogger(t *testing.T) {
	ec := Configuration{
		AppName:      "goProbe",
		HostName:     "localhost",
		HostPort:     5005,
		DebugEnabled: false,
		Interval:     120,
	}
	loadLogger(&ec)
	if goLog == nil {
		t.Errorf("Failed to load Logger, got : %t want : %t", false, true)
	}
}

func TestReadConfig(t *testing.T) {
	ec := Configuration{"goProbe", "localhost", 5005, "C:\\Logs\\error.log", false, 120}
	ac, al := ReadConfig()

	if ac.GetApplicationName() != ec.GetApplicationName() {
		t.Errorf("Failed to read configuration Application Name, got : %s want : %s", ac.GetApplicationName(), ec.GetApplicationName())
	}
	if ac.GetHostName() != ec.GetHostName() {
		t.Errorf("Failed to read configuration Host Name, got : %s want : %s", ac.GetHostName(), ec.GetHostName())
	}
	if ac.GetHostPort() != ec.GetHostPort() {
		t.Errorf("Failed to read configuration Host Port, got : %d want : %d", ac.GetHostPort(), ec.GetHostPort())
	}
	if ac.GetLogPath() != ec.GetLogPath() {
		t.Errorf("Failed to read configuration Log Path, got : %s want : %s", ac.GetLogPath(), ec.GetLogPath())
	}
	if ac.GetDebugEnabled() != ec.GetDebugEnabled() {
		t.Errorf("Failed to read configuration Debug Enabled, got : %t want : %t", ac.GetDebugEnabled(), ec.GetDebugEnabled())
	}
	if ac.GetInterval() != ec.GetInterval() {
		t.Errorf("Failed to read configuration Interval, got : %d want : %d", ac.GetInterval(), ec.GetInterval())
	}
	if al == nil {
		t.Errorf("Failed to load Logger, got : %t want : %t", false, true)
	}

}
