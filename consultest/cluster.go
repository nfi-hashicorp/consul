package consultest

type ConsulVersion string

const (
	V1_14_1 = "1.14.1"
)

type Cluster interface {
	// TODO
}

type BasicCluster struct {
	Replicas int
	Version  ConsulVersion
}

type EnterpriseCluster struct {
	BasicCluster
}
