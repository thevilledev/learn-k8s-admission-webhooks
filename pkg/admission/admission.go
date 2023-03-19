package admission

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/thevilledev/learn-admission-controllers/pkg/functions"
	"github.com/thevilledev/learn-admission-controllers/pkg/types"
	k8sadmission "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)

func Admit(body []byte) ([]byte, error) {
	var ar k8sadmission.AdmissionReview

	if _, _, err := universalDeserializer.Decode(body, nil, &ar); err != nil {
		return nil, fmt.Errorf("could not deserialize request: %v", err)
	} else if ar.Request == nil {
		return nil, errors.New("malformed admission review: request is nil")
	}

	arResp := k8sadmission.AdmissionReview{
		TypeMeta: ar.TypeMeta,
		Response: &k8sadmission.AdmissionResponse{
			UID: ar.Request.UID,
		},
	}

	var patchOps []types.PatchOperation
	for _, f := range functions.Registry {
		status, pops, err := f(ar.Request)
		if err != nil {
			arResp.Response.Allowed = false
			arResp.Response.Result = status
			break
		}
		if len(pops) > 0 {
			patchOps = append(patchOps, pops...)
		}
	}

	bytes, err := json.Marshal(&arResp)
	if err != nil {
		return nil, fmt.Errorf("marshaling response: %v", err)
	}
	return bytes, nil
}
