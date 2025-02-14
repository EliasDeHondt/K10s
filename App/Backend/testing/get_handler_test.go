package testing

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/handlers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var getClient = kubernetes.TestFakeClient()

func TestGetNodes(t *testing.T) {
	nodes, err := handlers.GetNodes(getClient)

	assert.NoError(t, err)
	assert.NotEmpty(t, nodes)
	assert.Equal(t, "node-1", (*nodes)[0].Name)
	assert.Equal(t, "node-2", (*nodes)[1].Name)
}

func TestGetPods(t *testing.T) {
	pods, err := handlers.GetPods(getClient, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, pods)
	assert.Equal(t, "pod-1", (*pods)[0].Name)
	assert.Equal(t, "pod-2", (*pods)[1].Name)
}

func TestGetPodsWithNamespace(t *testing.T) {
	pods, err := handlers.GetPods(getClient, "test")

	assert.NoError(t, err)
	assert.NotEmpty(t, pods)
	assert.Equal(t, "pod-3", (*pods)[0].Name)
}

func TestGetServices(t *testing.T) {
	services, err := handlers.GetServices(getClient, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, services)
	assert.Equal(t, "service-1", (*services)[0].Name)
}

func TestGetServicesWithNamespace(t *testing.T) {
	services, err := handlers.GetServices(getClient, "test")

	assert.NoError(t, err)
	assert.NotEmpty(t, services)
	assert.Equal(t, "service-2", (*services)[0].Name)
}

func TestGetDeployments(t *testing.T) {
	deployments, err := handlers.GetDeployments(getClient, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, deployments)
	assert.Equal(t, "deployment-1", (*deployments)[0].Name)
}

func TestGetDeploymentsWithNamespace(t *testing.T) {
	deployments, err := handlers.GetDeployments(getClient, "test")

	assert.NoError(t, err)
	assert.NotEmpty(t, deployments)
	assert.Equal(t, "deployment-2", (*deployments)[0].Name)
}

func TestGetConfigMaps(t *testing.T) {
	maps, err := handlers.GetConfigMaps(getClient, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, maps)
	assert.Equal(t, "configmap-1", (*maps)[0].Name)
}

func TestGetConfigMapsWithNamespace(t *testing.T) {
	maps, err := handlers.GetConfigMaps(getClient, "test")

	assert.NoError(t, err)
	assert.NotEmpty(t, maps)
	assert.Equal(t, "configmap-2", (*maps)[0].Name)
}

func TestGetSecrets(t *testing.T) {
	secrets, err := handlers.GetSecrets(getClient, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, secrets)
	assert.Equal(t, "secret-1", (*secrets)[0].Name)
}

func TestGetSecretsWithNamespace(t *testing.T) {
	secrets, err := handlers.GetSecrets(getClient, "test")

	assert.NoError(t, err)
	assert.NotEmpty(t, secrets)
	assert.Equal(t, "secret-2", (*secrets)[0].Name)
}

func TestGetTotalUsage(t *testing.T) {
	metrics, err := getClient.GetTotalUsage()

	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

func TestGetUsageForNode(t *testing.T) {
	metrics, err := getClient.GetUsageForNode("node-1")

	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

func TestGetUsageForNonExistingNode(t *testing.T) {
	metrics, err := getClient.GetUsageForNode("node-123")

	assert.Error(t, err)
	assert.Empty(t, metrics)
}
