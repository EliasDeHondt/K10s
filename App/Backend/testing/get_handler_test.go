package testing

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/handlers"
	"github.com/stretchr/testify/assert"
	"testing"
)

var client = kubernetes.TestFakeClient()

func TestGetNodes(t *testing.T) {
	nodes, err := handlers.GetNodes(&client, 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, nodes)
	assert.Equal(t, "node-1", (*nodes).Response[0].Name)
	assert.Equal(t, "node-2", (*nodes).Response[1].Name)
}

func TestGetPods(t *testing.T) {
	pods, err := handlers.GetPods(&client, "", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, pods)
	assert.Equal(t, "pod-1", (*pods).Response[0].Name)
	assert.Equal(t, "pod-2", (*pods).Response[1].Name)
}

func TestGetPodsWithNamespace(t *testing.T) {
	pods, err := handlers.GetPods(&client, "test", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, pods)
	assert.Equal(t, "pod-3", (*pods).Response[0].Name)
}

func TestGetServices(t *testing.T) {
	services, err := handlers.GetServices(&client, "", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, services)
	assert.Equal(t, "service-1", (*services).Response[0].Name)
}

func TestGetServicesWithNamespace(t *testing.T) {
	services, err := handlers.GetServices(&client, "test", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, services)
	assert.Equal(t, "service-2", (*services).Response[0].Name)
}

func TestGetDeployments(t *testing.T) {
	deployments, err := handlers.GetDeployments(&client, "", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, deployments)
	assert.Equal(t, "deployment-1", (*deployments).Response[0].Name)
}

func TestGetDeploymentsWithNamespace(t *testing.T) {
	deployments, err := handlers.GetDeployments(&client, "test", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, deployments)
	assert.Equal(t, "deployment-2", (*deployments).Response[0].Name)
}

func TestGetConfigMaps(t *testing.T) {
	maps, err := handlers.GetConfigMaps(&client, "", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, maps)
	assert.Equal(t, "configmap-1", (*maps).Response[0].Name)
}

func TestGetConfigMapsWithNamespace(t *testing.T) {
	maps, err := handlers.GetConfigMaps(&client, "test", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, maps)
	assert.Equal(t, "configmap-2", (*maps).Response[0].Name)
}

func TestGetSecrets(t *testing.T) {
	secrets, err := handlers.GetSecrets(&client, "", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, secrets)
	assert.Equal(t, "secret-1", (*secrets).Response[0].Name)
}

func TestGetSecretsWithNamespace(t *testing.T) {
	secrets, err := handlers.GetSecrets(&client, "test", 20, "")

	assert.NoError(t, err)
	assert.NotEmpty(t, secrets)
	assert.Equal(t, "secret-2", (*secrets).Response[0].Name)
}

func TestGetTotalUsage(t *testing.T) {
	metrics, err := client.GetTotalUsage()

	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

func TestGetUsageForNode(t *testing.T) {
	metrics, err := client.GetUsageForNode("node-1")

	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)
}

func TestGetUsageForNonExistingNode(t *testing.T) {
	metrics, err := client.GetUsageForNode("node-123")

	assert.Error(t, err)
	assert.Empty(t, metrics)
}
