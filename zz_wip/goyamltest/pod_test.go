package goyamltest

import (
	"encoding/json"
	"fmt"
	"testing"

	goccyyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
	k8syaml "sigs.k8s.io/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Test_pod(t *testing.T) {

	p := corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apiversion",
			Kind:       "kind",
		},
		ObjectMeta: metav1.ObjectMeta{
			CreationTimestamp: metav1.Time{},
			UID:               "xxx",
			Namespace:         "default",
			Name:              "xxx",
		},
		Spec: corev1.PodSpec{
			NodeSelector: map[string]string{
				"c": "d",
				"a": "b",
			},
		},
	}

	var v []byte
	v, _ = yamlv2.Marshal(p)
	fmt.Println("====yaml v2====")
	fmt.Println(string(v))

	v, _ = yamlv3.Marshal(p)
	fmt.Println("====yaml v3====")
	fmt.Println(string(v))

	v, _ = goccyyaml.Marshal(p)
	fmt.Println("====goccyyaml====")
	fmt.Println(string(v))

	v, _ = k8syaml.Marshal(p)
	fmt.Println("====k8syaml====")
	fmt.Println(string(v))

	pt := `
metadata:
  namespace: default
  name: xxx
apiVersion: "apiversion"
kind: "kind"
spec:
  containers:
  - name: test
	- name: test2
  nodeSelector:
    a: b
`

	p2 := &corev1.Pod{}
	goccyyaml.Unmarshal([]byte(pt), &p2)

	p2.Status = corev1.PodStatus{}

	v2, _ := k8syaml.Marshal(p2)
	fmt.Println(string(v2))

}

func TestYamlToUnstructured_Float64PortNumber(t *testing.T) {
	jsonStr := `{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "test-svc"
  },
  "spec": {
    "ports": [
      {
        "port": 8080.0,
        "protocol": "UDP"
      }
    ]
  }
}`

	var getPort = func(u *unstructured.Unstructured) interface{} {
		return u.Object["spec"].(map[string]interface{})["ports"].([]interface{})[0].(map[string]interface{})["port"]
	}

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	// Method 1
	// Summary, to unmarshal k8s `json or yaml string` to Unstructured, please unmarshal to Unstructured, not Unstructured.Object.
	if err := k8syaml.Unmarshal([]byte(jsonStr), &u); err != nil {
		t.Errorf("method 1 err: %s", err)
	}
	assert.Equal(t, int64(8080), getPort(u).(int64), "method 1 not equal")

	// Method 2
	if err := k8syaml.Unmarshal([]byte(jsonStr), &u.Object); err != nil {
		t.Errorf("method 2 err: %s", err)
	}
	assert.Equal(t, float64(8080), getPort(u).(float64), "method 2 not equal")

	// Method 3
	if err := json.Unmarshal([]byte(jsonStr), &u); err != nil {
		t.Errorf("method 3 err: %s", err)
	}
	assert.Equal(t, float64(8080), getPort(u).(float64), "method 3 not equal")

	// Method 4
	if err := json.Unmarshal([]byte(jsonStr), &u.Object); err != nil {
		t.Errorf("method 4 err: %s", err)
	}
	assert.Equal(t, float64(8080), getPort(u).(float64), "method 4 not equal")

}

func TestYamlToUnstructured_IntergerPortNumber(t *testing.T) {
	jsonStr := `{
  "apiVersion": "v1",
  "kind": "Service",
  "metadata": {
    "name": "test-svc"
  },
  "spec": {
    "ports": [
      {
        "port": 8080,
        "protocol": "UDP"
      }
    ]
  }
}`

	var getPort = func(u *unstructured.Unstructured) interface{} {
		return u.Object["spec"].(map[string]interface{})["ports"].([]interface{})[0].(map[string]interface{})["port"]
	}

	u := &unstructured.Unstructured{Object: map[string]interface{}{}}

	// Method 1
	// Summary, to unmarshal k8s `json or yaml string` to Unstructured, please unmarshal to Unstructured, not Unstructured.Object.
	if err := k8syaml.Unmarshal([]byte(jsonStr), &u); err != nil {
		t.Errorf("method 1 err: %s", err)
	}
	assert.Equal(t, int64(8080), getPort(u).(int64), "method 1 not equal")

	// Method 2
	if err := k8syaml.Unmarshal([]byte(jsonStr), &u.Object); err != nil {
		t.Errorf("method 2 err: %s", err)
	}
	assert.Equal(t, float64(8080), getPort(u).(float64), "method 2 not equal")

	// Method 3
	if err := json.Unmarshal([]byte(jsonStr), &u); err != nil {
		t.Errorf("method 3 err: %s", err)
	}
	assert.Equal(t, int64(8080), getPort(u).(int64), "method 3 not equal")

	// Method 4
	if err := json.Unmarshal([]byte(jsonStr), &u.Object); err != nil {
		t.Errorf("method 4 err: %s", err)
	}
	assert.Equal(t, float64(8080), getPort(u).(float64), "method 4 not equal")

}
