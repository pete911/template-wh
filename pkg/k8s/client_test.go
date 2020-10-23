package k8s

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestClient_GetConfigMap(t *testing.T) {

	t.Run("when returned config map is nil then no error is returned", func(t *testing.T) {

		configMapMock := new(ConfigMapsMock)
		configMapMock.On("Get", context.Background(), "vault-auth-roles", meta.GetOptions{}).Return(nil, nil)
		c := Client{configMapsGetter: &ConfigMapsGetterMock{getter: configMapMock}}

		cm, err := c.GetConfigMapData("kube-system", "vault-auth-roles")
		require.NoError(t, err)
		assert.Nil(t, cm)
	})

	t.Run("when get config map fails then error is returned", func(t *testing.T) {

		configMapMock := new(ConfigMapsMock)
		configMapMock.On("Get", context.Background(), "vault-auth-roles", meta.GetOptions{}).Return(nil, errors.New("test failuer"))
		c := Client{configMapsGetter: &ConfigMapsGetterMock{getter: configMapMock}}

		_, err := c.GetConfigMapData("kube-system", "vault-auth-roles")
		require.Error(t, err)
	})

	t.Run("when get config map is successful then no error is returned", func(t *testing.T) {

		expectedConfigMap := &v1.ConfigMap{
			Data: map[string]string{
				"test-role": `{"bound_service_account_names": ["default"], "bound_service_account_namespaces": ["*"], "token_policies": ["test"]}`,
			},
		}
		configMapMock := new(ConfigMapsMock)
		configMapMock.On("Get", context.Background(), "vault-auth-roles", meta.GetOptions{}).Return(expectedConfigMap, nil)
		c := Client{configMapsGetter: &ConfigMapsGetterMock{getter: configMapMock}}

		actualConfigMapData, err := c.GetConfigMapData("kube-system", "vault-auth-roles")
		require.NoError(t, err)
		assert.Equal(t, expectedConfigMap.Data, actualConfigMapData)
	})
}

// --- mocks ---

type ConfigMapsMock struct {
	mock.Mock
}

func (m ConfigMapsMock) Get(ctx context.Context, name string, options meta.GetOptions) (*v1.ConfigMap, error) {

	args := m.Called(ctx, name, options)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*v1.ConfigMap), args.Error(1)
}

type ConfigMapsGetterMock struct {
	getter *ConfigMapsMock
}

func (c ConfigMapsGetterMock) ConfigMaps(namespace string) configMapsInterface {
	return c.getter
}
