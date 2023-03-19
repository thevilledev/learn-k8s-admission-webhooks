package types

import (
	k8sadmission "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AdmitFunc func(req *k8sadmission.AdmissionRequest) (*metav1.Status, []PatchOperation, error)

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
