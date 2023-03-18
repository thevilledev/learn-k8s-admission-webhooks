package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	admission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
	podResource           = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
)

// patchOperation is an operation of a JSON patch, see https://tools.ietf.org/html/rfc6902 .
/*type patchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}*/

// isKubeNamespace checks if the given namespace is a Kubernetes-owned namespace.
func isKubeNamespace(ns string) bool {
	return ns == metav1.NamespacePublic || ns == metav1.NamespaceSystem
}

func Func(body []byte) ([]byte, error) {
	var rev admission.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &rev); err != nil {
		return nil, fmt.Errorf("could not deserialize request: %v", err)
	} else if rev.Request == nil {
		return nil, errors.New("malformed admission review: request is nil")
	}

	req := rev.Request
	if req.Resource != podResource {
		log.Printf("expect resource to be %s", podResource)
		return nil, errors.New("not a pod")
	}

	admissionReviewResponse := admission.AdmissionReview{
		// Since the admission webhook now supports multiple API versions, we need
		// to explicitly include the API version in the response.
		// This API version needs to match the version from the request exactly, otherwise
		// the API server will be unable to process the response.
		// Note: a v1beta1 AdmissionReview is JSON-compatible with the v1 version, that's why
		// we do not need to differentiate during unmarshaling or in the actual logic.
		TypeMeta: rev.TypeMeta,
		Response: &admission.AdmissionResponse{
			UID: rev.Request.UID,
		},
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	if _, ok := pod.Labels["mandatory"]; !ok {
		admissionReviewResponse.Response.Allowed = false
		admissionReviewResponse.Response.Result = &metav1.Status{
			Message: "missing required label",
		}
	}

	/*var patchOps []patchOperation
	// Apply the admit() function only for non-Kubernetes namespaces. For objects in Kubernetes namespaces, return
	// an empty set of patch operations.
	if !isKubeNamespace(admissionReviewReq.Request.Namespace) {
		patchOps, err = patchFunc(admissionReviewReq.Request)
	}*/

	//var patchOps []patchOperation
	// Apply the admit() function only for non-Kubernetes namespaces. For objects in Kubernetes namespaces, return
	// an empty set of patch operations.
	if isKubeNamespace(rev.Request.Namespace) {
		//patchOps, err := patchFunc(admissionReviewReq.Request)
		admissionReviewResponse.Response.Allowed = true
	}

	// Return the AdmissionReview with a response as JSON.
	bytes, err := json.Marshal(&admissionReviewResponse)
	if err != nil {
		return nil, fmt.Errorf("marshaling response: %v", err)
	}
	return bytes, nil
}
