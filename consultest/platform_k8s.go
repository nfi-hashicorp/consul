package consultest

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/hashicorp/consul/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// implements Platform
type K8sPlatform struct{}

func (p *K8sPlatform) Supports(c Cluster) bool {
	panic("not implemented") // TODO: Implement
}

func (p *K8sPlatform) Deploy(c Cluster) Deployment {
	panic("not implemented") // TODO: Implement
}

func (p *K8sPlatform) Cleanup() {
	panic("not implemented") // TODO: Implement
}

// implements Deployment
type CTIAK8sDeployment struct {
	Namespace    string
	k8sClientSet *kubernetes.Clientset
	client       *api.Client
}

func (d *CTIAK8sDeployment) Client() *api.Client {
	panic("nyi")
}

func (d *CTIAK8sDeployment) Cleanup() {
	// TODO: delete client (and namespace?)
}

// CTIAK8sPlatform targets the CTIA Kubernetes baseline platform.
//
// Implements Platform
type CTIAK8sPlatform struct {
	deployments []CTIAK8sDeployment
}

func (p *CTIAK8sPlatform) k8sClientSet() *kubernetes.Clientset {
	// TODO: use AWS SDK or something
	home, _ := os.UserHomeDir()
	ctiaK8sConfigPath := filepath.Join(home, ".kube/ctia-shared.kubeconfig")
	config, err := clientcmd.BuildConfigFromFlags("", ctiaK8sConfigPath)
	if err != nil {
		log.Fatalf("loading kubeconfig: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("creating k8s clientset: %s", err)
	}
	return clientset
}

func (p *CTIAK8sPlatform) Deploy(c Cluster) Deployment {
	if p.deployments == nil {
		p.deployments = []CTIAK8sDeployment{}
	}

	// TODO: a little weird
	clientset := p.k8sClientSet()

	ns := fmt.Sprintf("test-%s", strings.ToLower(random.UniqueId()))
	_, err := clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("failed to create namespace %q: %s", ns, err)
	}

	d := CTIAK8sDeployment{
		Namespace:    ns,
		k8sClientSet: clientset,
	}
	p.deployments = append(p.deployments, d)
	// TODO: create tunnel and consul client
	// TODO: deploy
	return &d
}

func (p *CTIAK8sPlatform) Cleanup() {
	for len(p.deployments) > 0 {
		d := p.deployments[0]
		p.deployments = p.deployments[1:]
		d.k8sClientSet.CoreV1().Namespaces().Delete(context.Background(), d.Namespace, metav1.DeleteOptions{})
	}
}

func (p *CTIAK8sPlatform) Supports(c Cluster) bool {
	switch c.(type) {
	case *BasicCluster:
		return true
	}
	return false
}
