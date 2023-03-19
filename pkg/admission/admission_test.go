package admission

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	k8sadmission "k8s.io/api/admission/v1"
)

func TestValidMutationRequest(t *testing.T) {
	payload := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
			"uid": "7f0b2891-916f-4ed6-b7cd-27bff1815a8c",
			"kind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"resource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"requestKind": {
				"group": "",
				"version": "v1",
				"kind": "Pod"
			},
			"requestResource": {
				"group": "",
				"version": "v1",
				"resource": "pods"
			},
			"namespace": "yolo",
			"operation": "CREATE",
			"userInfo": {
				"username": "kubernetes-admin",
				"groups": [
					"system:masters",
					"system:authenticated"
				]
			},
			"object": {
				"kind": "Pod",
				"apiVersion": "v1",
				"metadata": {
					"name": "c7m",
					"namespace": "yolo",
					"creationTimestamp": null,
					"labels": {
						"name": "c7m"
					},
					"annotations": {
						"kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Pod\",\"metadata\":{\"annotations\":{},\"labels\":{\"name\":\"c7m\"},\"name\":\"c7m\",\"namespace\":\"yolo\"},\"spec\":{\"containers\":[{\"args\":[\"-c\",\"trap \\\"killall sleep\\\" TERM; trap \\\"kill -9 sleep\\\" KILL; sleep infinity\"],\"command\":[\"/bin/bash\"],\"image\":\"centos:7\",\"name\":\"c7m\"}]}}\n"
					}
				},
				"spec": {
					"volumes": [
						{
							"name": "default-token-5z7xl",
							"secret": {
								"secretName": "default-token-5z7xl"
							}
						}
					],
					"containers": [
						{
							"name": "c7m",
							"image": "centos:7",
							"command": [
								"/bin/bash"
							],
							"args": [
								"-c",
								"trap \"killall sleep\" TERM; trap \"kill -9 sleep\" KILL; sleep infinity"
							],
							"resources": {},
							"volumeMounts": [
								{
									"name": "default-token-5z7xl",
									"readOnly": true,
									"mountPath": "/var/run/secrets/kubernetes.io/serviceaccount"
								}
							],
							"terminationMessagePath": "/dev/termination-log",
							"terminationMessagePolicy": "File",
							"imagePullPolicy": "IfNotPresent"
						}
					],
					"restartPolicy": "Always",
					"terminationGracePeriodSeconds": 30,
					"dnsPolicy": "ClusterFirst",
					"serviceAccountName": "default",
					"serviceAccount": "default",
					"securityContext": {},
					"schedulerName": "default-scheduler",
					"tolerations": [
						{
							"key": "node.kubernetes.io/not-ready",
							"operator": "Exists",
							"effect": "NoExecute",
							"tolerationSeconds": 300
						},
						{
							"key": "node.kubernetes.io/unreachable",
							"operator": "Exists",
							"effect": "NoExecute",
							"tolerationSeconds": 300
						}
					],
					"priority": 0,
					"enableServiceLinks": true
				},
				"status": {}
			},
			"oldObject": null,
			"dryRun": false,
			"options": {
				"kind": "CreateOptions",
				"apiVersion": "meta.k8s.io/v1"
			}
		}
	}`
	resp, err := Admit([]byte(payload))
	if err != nil {
		t.Errorf("AdmissionRequest failed with error: %s", err)
	}

	r := k8sadmission.AdmissionReview{}
	err = json.Unmarshal(resp, &r)
	assert.NoError(t, err, "failed to unmarshal with error: %s", err)

	rr := r.Response
	assert.Equal(t, `[{"op":"replace","path":"/spec/containers/0/image","value":"debian:latest"}]`, string(rr.Patch))
	assert.Equal(t, true, rr.Allowed)
}

func TestErrorsOnInvalidJson(t *testing.T) {
	payload := `Everything here is so cold / Everything here is so dark`
	_, err := Admit([]byte(payload))
	if err == nil {
		t.Error("did not fail when sending invalid json")
	}
}

func TestEmptyRequest(t *testing.T) {
	payload := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {}
	}`
	resp, err := Admit([]byte(payload))
	assert.NoError(t, err)

	r := k8sadmission.AdmissionReview{}
	err = json.Unmarshal(resp, &r)
	assert.NoError(t, err, "failed to unmarshal with error: %s", err)

	assert.Equal(t, true, r.Response.Allowed)
}
