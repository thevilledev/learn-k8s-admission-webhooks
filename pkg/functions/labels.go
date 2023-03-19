package functions

import (
	"fmt"
	"log"

	"github.com/thevilledev/learn-admission-controllers/pkg/types"
	k8sadmission "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func init() {
	registerFunction("verify pod labels", verifyPodLabels)
}

var verifyPodLabels types.AdmitFunc = func(req *k8sadmission.AdmissionRequest) (*metav1.Status, []types.PatchOperation, error) {
	if req.Resource != podResource {
		log.Printf("Ignore admission request %s as it's not a pod resource", string(req.UID))
		return nil, nil, nil
	}

	r := req.Object.Raw
	p := corev1.Pod{}
	if _, _, err := universalDeserializer.Decode(r, nil, &p); err != nil {
		return nil, nil, fmt.Errorf("could not deserialize pod object: %v", err)
	}

	if _, found := p.Labels["foo"]; !found {
		return &metav1.Status{Message: "missing required label"}, nil, fmt.Errorf("could not find label 'foo' from pod")
	}

	return nil, nil, nil
}
