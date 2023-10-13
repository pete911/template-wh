package k8s

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattbaird/jsonpatch"
	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log/slog"
	"os"
	"strings"
)

// Manifest is needed, so we don't replace values in metadata["managedFields"] and so on
type Manifest struct {
	Metadata Metadata        `json:"metadata,omitempty"`
	Spec     json.RawMessage `json:"spec,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
}

type Metadata struct {
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	ClusterName string            `json:"clusterName,omitempty"`
}

// Mutate takes AdmissionReview and replaces AdmissionReview.Request $var with supplied values
// (although this can replace ${var} as well, it won't play nicely with json payload)
func Mutate(body []byte, values map[string]string) ([]byte, error) {

	var admissionReview admission.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("mutate request unmarshal: %w", err)
	}

	if admissionReview.Request == nil {
		return nil, errors.New("received nil admission review request")
	}

	var manifest Manifest
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, &manifest); err != nil {
		return nil, fmt.Errorf("unmrshal manifest: %v", err)
	}

	patch, err := getPatch(manifest, values)
	if err != nil {
		return nil, fmt.Errorf("get patch: %v", err)
	}

	admissionReview.Response = getAdmissionResponse(admissionReview.Request.UID, patch)
	responseBody, err := json.Marshal(admissionReview)
	if err != nil {
		return nil, fmt.Errorf("mutate response marshal: %w", err)
	}
	return responseBody, nil
}

func getPatch(manifest Manifest, values map[string]string) ([]byte, error) {

	oldObj, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("marshal manifest: %v", err)
	}
	newObj, err := json.Marshal(manifest)
	if err != nil {
		return nil, fmt.Errorf("marshal manifest: %v", err)
	}

	newObj = expand(newObj, values)
	return createPatch(oldObj, newObj)
}

// Expand replaces ${var} or $var in the request with passed values
func expand(request []byte, values map[string]string) []byte {

	// TODO - value can contain characters that would break json e.g. {, [
	return []byte(os.Expand(string(request), func(key string) string {
		return strings.TrimSpace(values[key])
	}))
}

// Create json patch
// https://kubernetes.io/blog/2019/03/21/a-guide-to-kubernetes-admission-controllers/#object-modification-logic
func createPatch(a, b []byte) ([]byte, error) {

	patch, err := jsonpatch.CreatePatch(a, b)
	if err != nil {
		return nil, fmt.Errorf("create patch: %w", err)
	}
	if len(patch) == 0 {
		return nil, nil
	}
	rawPatch, err := json.Marshal(patch)
	if err != nil {
		return nil, fmt.Errorf("marshal %+v patch: %w", patch, err)
	}
	return rawPatch, nil
}

func getAdmissionResponse(admissionRequestUID types.UID, patch []byte) *admission.AdmissionResponse {

	patchTypeJson := admission.PatchTypeJSONPatch
	var admissionResponse = &admission.AdmissionResponse{
		Allowed: true,
		UID:     admissionRequestUID,
		Patch:   patch,
		Result: &meta.Status{
			Status: "Success",
		},
	}

	// patch type can be set only if there is actual patch
	if patch != nil {
		slog.Info(fmt.Sprintf("patch: %+v", string(patch)))
		admissionResponse.AuditAnnotations = map[string]string{"mutated": "template-wh"}
		admissionResponse.PatchType = &patchTypeJson
	}
	return admissionResponse
}
