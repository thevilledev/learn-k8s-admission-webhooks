package functions

import (
	"fmt"
	"log"

	"github.com/thevilledev/learn-admission-controllers/pkg/types"
	k8sadmission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	greatestImage = "debian:latest"
)

func init() {
	registerFunction("set image", setImage)
}

var setImage types.AdmitFunc = func(req *k8sadmission.AdmissionRequest) (*metav1.Status, []types.PatchOperation, error) {
	if req.Resource != podResource {
		log.Printf("Ignore admission request %s as it's not a pod resource", string(req.UID))
		return nil, nil, nil
	}

	r := req.Object.Raw
	p := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(r, nil, &p); err != nil {
		return nil, nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	var patches []types.PatchOperation
	if p.Spec.Containers[0].Image != greatestImage {
		patches = append(patches, types.PatchOperation{
			Op:    "replace",
			Path:  "/spec/containers/0/image",
			Value: greatestImage,
		})
	}

	return nil, patches, nil
}
