package kubernetes

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/kubermatic/api"
	"github.com/kubermatic/api/provider"
	kapi "k8s.io/kubernetes/pkg/api"
	kerrors "k8s.io/kubernetes/pkg/api/errors"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/labels"
)

var _ provider.KubernetesProvider = (*kubernetesProvider)(nil)

type kubernetesProvider struct {
	mu     sync.Mutex
	cps    map[string]provider.CloudProvider
	client *client.Client
}

// NewKubernetesProvider creates a new kubernetes provider object
func NewKubernetesProvider(
	clientConfig *client.Config,
	cps map[string]provider.CloudProvider,
) provider.KubernetesProvider {
	client, err := client.New(clientConfig)
	if err != nil {
		log.Fatal(err)
	}

	return &kubernetesProvider{
		cps:    cps,
		client: client,
	}
}

func (p *kubernetesProvider) NewCluster(user, cluster string, spec *api.ClusterSpec) (*api.Cluster, error) {
	// call cluster before lock is taken
	cs, err := p.Clusters(user)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// sanity checks for a fresh cluster
	switch {
	case cluster == "":
		return nil, kerrors.NewBadRequest("cluster name is required")
	case user == "":
		return nil, kerrors.NewBadRequest("cluster user is required")
	case spec.HumanReadableName == "":
		return nil, kerrors.NewBadRequest("cluster humanReadableName is required")
	}

	for _, c := range cs {
		if c.Spec.HumanReadableName == spec.HumanReadableName {
			return nil, kerrors.NewAlreadyExists("cluster", spec.HumanReadableName)
		}
	}

	ns := &kapi.Namespace{
		ObjectMeta: kapi.ObjectMeta{
			Name:        NamespaceName(user, cluster),
			Annotations: map[string]string{},
			Labels:      map[string]string{},
		},
	}

	c := &api.Cluster{
		Metadata: api.Metadata{
			User: user,
			Name: cluster,
		},
		Spec: *spec,
		Status: api.ClusterStatus{
			LastTransitionTime: time.Now(),
			Phase:              api.PendingClusterStatusPhase,
		},
	}

	ns, err = MarshalCluster(p.cps, c, ns)
	if err != nil {
		return nil, err
	}

	ns, err = p.client.Namespaces().Create(ns)
	if err != nil {
		return nil, err
	}

	c, err = UnmarshalCluster(p.cps, ns)
	if err != nil {
		_ = p.client.Namespaces().Delete(NamespaceName(user, cluster))
		return nil, err
	}

	return c, nil
}

func (p *kubernetesProvider) Cluster(user, cluster string) (*api.Cluster, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	ns, err := p.client.Namespaces().Get(NamespaceName(user, cluster))
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil, kerrors.NewNotFound("cluster", cluster)
		}
		return nil, err
	}

	c, err := UnmarshalCluster(p.cps, ns)
	if err != nil {
		return nil, err
	}

	if c.Metadata.User != user {
		// don't return Forbidden, not NotFound to obfuscate the existence
		return nil, kerrors.NewNotFound("cluster", cluster)
	}

	return c, nil
}

func (p *kubernetesProvider) Clusters(user string) ([]*api.Cluster, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	nsList, err := p.client.Namespaces().List(
		labels.SelectorFromSet(labels.Set(map[string]string{
			RoleLabelKey: ClusterRoleLabel,
			userLabelKey: LabelUser(user),
		})),
		fields.Everything(),
	)
	if err != nil {
		return nil, err
	}

	cs := make([]*api.Cluster, 0, len(nsList.Items))
	for i := range nsList.Items {
		ns := nsList.Items[i]
		c, err := UnmarshalCluster(p.cps, &ns)
		if err != nil {
			log.Println(fmt.Sprintf("error unmarshaling namespace %s: %v", ns.Name, err))
			continue
		}

		cs = append(cs, c)
	}

	return cs, nil
}

func (p *kubernetesProvider) DeleteCluster(user, cluster string) error {
	// check permission by getting the cluster first
	_, err := p.Cluster(user, cluster)
	if err != nil {
		return err
	}

	return p.client.Namespaces().Delete(NamespaceName(user, cluster))
}

func (p *kubernetesProvider) Nodes(user, cluster string) ([]string, error) {
	return []string{}, nil
}
