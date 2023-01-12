package consultest

import (
	"fmt"
	"testing"

	"github.com/hashicorp/consul/api"
)

// Platforms know how to deploy (some types of) Clusters.
//
// Cleanup should be called before test end, e.g, with testing.T#Cleanup. Our TestMain
// does this automatically.
type Platform interface {
	Supports(Cluster) bool
	Deploy(Cluster) Deployment
	Cleanup()
}

type Deployment interface {
	Client() *api.Client
}

// TestMainPlatforms should be configured in callers' main_test.go#TestMain(m *testing.M).
// See TestMain for a convenience function
//
// After running tests, Platform.Cleanup() should be called on every platform.
var TestMPlatforms []Platform

// TestM can be called from a func TestMain(m *testing.M) (see https://pkg.go.dev/testing#hdr-Main)
// to set up TestMainPlatforms. It is responsible for:
func TestM(m *testing.M, platforms []Platform) int {
	TestMPlatforms = platforms
	code := m.Run()
	for _, p := range platforms {
		p.Cleanup()
	}
	return code
}

// Deploy deploys Cluster c to the first applicable platform passed to TestMain
func Deploy(c Cluster) Deployment {
	for _, p := range TestMPlatforms {
		if p.Supports(c) {
			return p.Deploy(c)
		}
	}
	panic(fmt.Sprintf("no platform supports Cluster %v", c))
}
