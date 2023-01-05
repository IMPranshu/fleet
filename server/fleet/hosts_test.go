package fleet

import (
	"fmt"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHostStatus(t *testing.T) {
	mockClock := clock.NewMockClock()

	testCases := []struct {
		seenTime            time.Time
		distributedInterval uint
		configTLSRefresh    uint
		status              HostStatus
	}{
		{mockClock.Now().Add(-30 * time.Second), 10, 3600, StatusOnline},
		{mockClock.Now().Add(-75 * time.Second), 10, 3600, StatusOffline},
		{mockClock.Now().Add(-30 * time.Second), 3600, 10, StatusOnline},
		{mockClock.Now().Add(-75 * time.Second), 3600, 10, StatusOffline},

		{mockClock.Now().Add(-60 * time.Second), 60, 60, StatusOnline},
		{mockClock.Now().Add(-121 * time.Second), 60, 60, StatusOffline},

		{mockClock.Now().Add(-1 * time.Second), 10, 10, StatusOnline},
		{mockClock.Now().Add(-2 * time.Minute), 10, 10, StatusOffline},
		{mockClock.Now().Add(-31 * 24 * time.Hour), 10, 10, StatusOffline}, // As of Fleet 4.15, StatusMIA is deprecated in favor of StatusOffline

		// Ensure behavior is reasonable if we don't have the values
		{mockClock.Now().Add(-1 * time.Second), 0, 0, StatusOnline},
		{mockClock.Now().Add(-2 * time.Minute), 0, 0, StatusOffline},
		{mockClock.Now().Add(-31 * 24 * time.Hour), 0, 0, StatusOffline}, // As of Fleet 4.15, StatusMIA is deprecated in favor of StatusOffline
	}

	for _, tt := range testCases {
		t.Run("", func(t *testing.T) {
			// Save interval values
			h := Host{
				DistributedInterval: tt.distributedInterval,
				ConfigTLSRefresh:    tt.configTLSRefresh,
				SeenTime:            tt.seenTime,
			}

			assert.Equal(t, tt.status, h.Status(mockClock.Now()))
		})
	}
}

func TestHostIsNew(t *testing.T) {
	mockClock := clock.NewMockClock()

	host := Host{}

	host.CreatedAt = mockClock.Now().AddDate(0, 0, -1)
	assert.True(t, host.IsNew(mockClock.Now()))

	host.CreatedAt = mockClock.Now().AddDate(0, 0, -2)
	assert.False(t, host.IsNew(mockClock.Now()))
}

func TestPlatformFromHost(t *testing.T) {
	for _, tc := range []struct {
		host        string
		expPlatform string
	}{
		{
			host:        "unknown",
			expPlatform: "",
		},
		{
			host:        "",
			expPlatform: "",
		},
		{
			host:        "linux",
			expPlatform: "linux",
		},
		{
			host:        "ubuntu",
			expPlatform: "linux",
		},
		{
			host:        "debian",
			expPlatform: "linux",
		},
		{
			host:        "rhel",
			expPlatform: "linux",
		},
		{
			host:        "centos",
			expPlatform: "linux",
		},
		{
			host:        "sles",
			expPlatform: "linux",
		},
		{
			host:        "kali",
			expPlatform: "linux",
		},
		{
			host:        "gentoo",
			expPlatform: "linux",
		},
		{
			host:        "darwin",
			expPlatform: "darwin",
		},
		{
			host:        "windows",
			expPlatform: "windows",
		},
	} {
		fleetPlatform := PlatformFromHost(tc.host)
		require.Equal(t, tc.expPlatform, fleetPlatform)

	}
}

func TestHostDisplayName(t *testing.T) {
	const (
		computerName   = "K0mpu73rN4M3"
		hostname       = "h0s7N4ME"
		hardwareModel  = "M0D3l"
		hardwareSerial = "53r14l"
	)
	for _, tc := range []struct {
		host     Host
		expected string
	}{
		{
			host:     Host{ComputerName: computerName, Hostname: hostname, HardwareModel: hardwareModel, HardwareSerial: hardwareSerial},
			expected: computerName, // If ComputerName is present, DisplayName is ComputerName
		},
		{
			host:     Host{ComputerName: "", Hostname: "h0s7N4ME", HardwareModel: "M0D3l", HardwareSerial: "53r14l"},
			expected: hostname, // If ComputerName is empty, DisplayName is Hostname (if present)
		},
		{
			host:     Host{ComputerName: "", Hostname: "", HardwareModel: "M0D3l", HardwareSerial: "53r14l"},
			expected: fmt.Sprintf("%s (%s)", hardwareModel, hardwareSerial), // If ComputerName and Hostname are empty, DisplayName is composite of HardwareModel and HardwareSerial (if both are present)
		},
		{
			host:     Host{ComputerName: "", Hostname: "", HardwareModel: "", HardwareSerial: hardwareSerial},
			expected: "", // If HarwareModel and/or HardwareSerial are empty, DisplayName is also empty
		},
		{
			host:     Host{ComputerName: "", Hostname: "", HardwareModel: hardwareModel, HardwareSerial: ""},
			expected: "", // If HarwareModel and/or HardwareSerial are empty, DisplayName is also empty
		},
		{
			host:     Host{ComputerName: "", Hostname: "", HardwareModel: "", HardwareSerial: ""},
			expected: "", // If HarwareModel and/or HardwareSerial are empty, DisplayName is also empty
		},
	} {
		require.Equal(t, tc.expected, tc.host.DisplayName())
	}
}

func TestMDMEnrollmentStatus(t *testing.T) {
	for _, tc := range []struct {
		hostMDM  HostMDM
		expected string
	}{
		{
			hostMDM:  HostMDM{Enrolled: true, InstalledFromDep: true},
			expected: "On (automatic)",
		},
		{
			hostMDM:  HostMDM{Enrolled: true, InstalledFromDep: false},
			expected: "On (manual)",
		},
		{
			hostMDM:  HostMDM{Enrolled: false, InstalledFromDep: true},
			expected: "Pending",
		},
		{
			hostMDM:  HostMDM{Enrolled: false, InstalledFromDep: false},
			expected: "Off",
		},
	} {
		require.Equal(t, tc.expected, tc.hostMDM.EnrollmentStatus())
	}
}
