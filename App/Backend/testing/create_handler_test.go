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

func TestCreateNode(t *testing.T) {
	data := getData("Node.yaml")

	resources, err := handlers.CreateResources(createClient, data)

	assert.NoError(t, err)
	assert.NotNil(t, resources)
	assert.Equal(t, 1, len(resources))
}

func TestCreatePod(t *testing.T) {

}
