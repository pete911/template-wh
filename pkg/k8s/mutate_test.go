package k8s

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	admission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

func TestMutate(t *testing.T) {

	t.Run("when request contains placeholders then patch is created", func(t *testing.T) {

		b, err := Mutate(getAdmissionReview(t, helmReleaseWithPlaceholder), values)
		require.NoError(t, err)

		var response admission.AdmissionReview
		require.NoError(t, json.Unmarshal(b, &response))

		expectedPatch := `[{"op":"replace","path":"/spec/chart/repository","value":"https://kubernetes-charts.storage.googleapis.com/"}]`
		actualPatch := string(response.Response.Patch)
		assert.Equal(t, expectedPatch, actualPatch)
		assert.Equal(t, "template-wh", response.Response.AuditAnnotations["mutated"])
	})

	t.Run("when request does not contain placeholders then no patch is created", func(t *testing.T) {

		b, err := Mutate(getAdmissionReview(t, helmRelease), values)
		require.NoError(t, err)

		var response admission.AdmissionReview
		require.NoError(t, json.Unmarshal(b, &response))

		expectedPatch := ``
		actualPatch := string(response.Response.Patch)
		assert.Equal(t, expectedPatch, actualPatch)
		assert.Equal(t, 0, len(response.Response.AuditAnnotations))
	})
}

// --- test helpers ---

func getAdmissionReview(t *testing.T, object string) []byte {

	data := admission.AdmissionReview{
		Request: &admission.AdmissionRequest{
			UID: "abc123",
			Object: runtime.RawExtension{
				Raw: []byte(object),
			},
		},
	}

	out, err := json.Marshal(data)
	require.NoError(t, err)
	return out
}

// --- test data ---

var (
	values      = map[string]string{"host": "https://kubernetes-charts.storage.googleapis.com/"}
	helmRelease = `{
  "apiVersion": "helm.fluxcd.io/v1",
  "kind": "HelmRelease",
  "metadata": {
    "name": "rabbit",
    "namespace": "default"
  },
  "spec": {
    "releaseName": "rabbitmq",
    "targetNamespace": "mq",
    "timeout": 300,
    "resetValues": false,
    "wait": false,
    "forceUpgrade": false,
    "chart": {
      "repository": "https://kubernetes-charts.storage.googleapis.com/",
      "name": "rabbitmq",
      "version": "3.3.6"
    },
    "values": {
      "replicas": 1
    }
  }
}
`
	helmReleaseWithPlaceholder = `{
  "apiVersion": "helm.fluxcd.io/v1",
  "kind": "HelmRelease",
  "metadata": {
    "name": "rabbit",
    "namespace": "default"
  },
  "spec": {
    "releaseName": "rabbitmq",
    "targetNamespace": "mq",
    "timeout": 300,
    "resetValues": false,
    "wait": false,
    "forceUpgrade": false,
    "chart": {
      "repository": "$host",
      "name": "rabbitmq",
      "version": "3.3.6"
    },
    "values": {
      "replicas": 1
    }
  }
}
`
)
