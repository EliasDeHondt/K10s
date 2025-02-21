/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/
package testing

import (
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes"
	"github.com/eliasdehondt/K10s/App/Backend/cmd/kubernetes/handlers"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var createClient = kubernetes.TestFakeClient()

func getData(fileName string) []byte {
	data, err := os.ReadFile("test_yamls/" + fileName)
	if err != nil {
		panic(err)
	}
	return data
}

func TestCreateNamespace(t *testing.T) {
	data := getData("Namespace.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCreateNode(t *testing.T) {
	data := getData("Node.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCreatePod(t *testing.T) {
	data := getData("Pod.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCreateService(t *testing.T) {
	data := getData("Service.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCreateDeployment(t *testing.T) {
	data := getData("Deployment.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCreateConfigMap(t *testing.T) {
	data := getData("ConfigMap.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCreateSecret(t *testing.T) {
	data := getData("Secret.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(resources))
}

func TestCombined(t *testing.T) {
	data := getData("Combined.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.Equal(t, 5, len(resources))
}
