package injector

import (
	"net/http/httptest"
	"strings"
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/golang/mock/gomock"
	"github.com/openservicemesh/osm/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func BenchmarkPodCreationHandler(b *testing.B) {
	b.StopTimer()
	admissionRequestBody := `{
		"kind": "AdmissionReview",
		"apiVersion": "admission.k8s.io/v1",
		"request": {
		  "uid": "11111111-2222-3333-4444-555555555555",
		  "kind": {
			"group": "",
			"version": "v1",
			"kind": "PodExecOptions"
		  },
		  "resource": {
			"group": "",
			"version": "v1",
			"resource": "pods"
		  },
		  "subResource": "exec",
		  "requestKind": {
			"group": "",
			"version": "v1",
			"kind": "PodExecOptions"
		  },
		  "requestResource": {
			"group": "",
			"version": "v1",
			"resource": "pods"
		  },
		  "requestSubResource": "exec",
		  "name": "some-pod-1111111111-22222",
		  "namespace": "default",
		  "operation": "CONNECT",
		  "userInfo": {
			"username": "user",
			"groups": []
		  },
		  "object": {
			"kind": "PodExecOptions",
			"apiVersion": "v1",
			"stdin": true,
			"stdout": true,
			"tty": true,
			"container": "some-pod",
			"command": ["bin/bash"]
		  },
		  "oldObject": null,
		  "dryRun": false,
		  "options": null
		}
	  }`

	req := httptest.NewRequest("GET", "/a/b/c", strings.NewReader(admissionRequestBody))
	req.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	client := fake.NewSimpleClientset()
	mockNsController := k8s.NewMockController(gomock.NewController(b))
	mockNsController.EXPECT().GetNamespace("default").Return(&corev1.Namespace{})
	mockNsController.EXPECT().IsMonitoredNamespace("default").Return(true).Times(1)

	wh := &mutatingWebhook{
		kubeClient:          client,
		kubeController:      mockNsController,
		nonInjectNamespaces: mapset.NewSet(),
	}
	w := httptest.NewRecorder()

	b.StartTimer()
	wh.podCreationHandler(w, req)
	b.StopTimer()

	res := w.Result()
	if res.StatusCode != 200 {
		b.Errorf("Expected 200, got %d", res.StatusCode)
	}
}