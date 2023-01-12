package consultest

// DockerPlatform deploys to Docker on localhost
//
// Implements Platform
type DockerPlatform struct{}

func (p *DockerPlatform) Supports(c Cluster) bool {
	panic("not implemented") // TODO: Implement
}

func (p *DockerPlatform) Deploy(c Cluster) Deployment {
	panic("not implemented") // TODO: Implement
}

func (p *DockerPlatform) Cleanup() {
	panic("not implemented") // TODO: Implement
}
