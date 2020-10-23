package k8s

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mattbaird/jsonpatch"
	admission "k8s.io/api/admission/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
)

// Takes AdmissionReview and replaces AdmissionReview.Request $var with supplied values
// (although this can replace ${var} as well, it won't play nicely with json payload)
func Mutate(body []byte, values map[string]string) ([]byte, error) {

	var admissionReview admission.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		return nil, fmt.Errorf("mutate request unmarshal: %w", err)
	}

	if admissionReview.Request == nil {
		return nil, errors.New("received nil admission review request")
	}

	modified := expand(admissionReview.Request.Object.Raw, values)
	patch, err := createPatch(admissionReview.Request.Object.Raw, modified)
	if err != nil {
		return nil, err
	}

	admissionReview.Response = getAdmissionResponse(admissionReview.Request.UID, patch)
	responseBody, err := json.Marshal(admissionReview)
	if err != nil {
		return nil, fmt.Errorf("mutate response marshal: %w", err)
	}
	return responseBody, nil
}

// Expand replaces ${var} or $var in the request with passed values
func expand(request []byte, values map[string]string) []byte {

	return []byte(os.Expand(string(request), func(key string) string {
		return values[key]
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
		return nil, fmt.Errorf("marshal patch: %w", err)
	}
	return rawPatch, nil
}

func getAdmissionResponse(admissionRequestUID types.UID, patch []byte) *admission.AdmissionResponse {

	pt := admission.PatchTypeJSONPatch
	var auditAnnotations map[string]string
	if patch != nil {
		auditAnnotations = map[string]string{"mutated": "template-wh"}
	}

	return &admission.AdmissionResponse{
		Allowed:          true,
		UID:              admissionRequestUID,
		PatchType:        &pt,
		AuditAnnotations: auditAnnotations,
		Patch:            patch,
		Result: &meta.Status{
			Status: "Success",
		},
	}
}
