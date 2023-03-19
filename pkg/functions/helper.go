package functions

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	podResource           = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
	universalDeserializer = serializer.NewCodecFactory(runtime.NewScheme()).UniversalDeserializer()
)
