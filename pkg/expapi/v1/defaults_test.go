package v1

import (
	"reflect"
	"testing"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/runtime"
)

func TestSetDefaultDaemon(t *testing.T) {
	tests := []struct {
		dc                 *Daemon
		expectLabelsChange bool
	}{
		{
			dc: &Daemon{
				Spec: DaemonSpec{
					Template: &v1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectLabelsChange: true,
		},
		{
			dc: &Daemon{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"bar": "foo",
					},
				},
				Spec: DaemonSpec{
					Template: &v1.PodTemplateSpec{
						ObjectMeta: v1.ObjectMeta{
							Labels: map[string]string{
								"foo": "bar",
							},
						},
					},
				},
			},
			expectLabelsChange: false,
		},
	}

	for _, test := range tests {
		dc := test.dc
		obj2 := roundTrip(t, runtime.Object(dc))
		dc2, ok := obj2.(*Daemon)
		if !ok {
			t.Errorf("unexpected object: %v", dc2)
			t.FailNow()
		}
		if test.expectLabelsChange != reflect.DeepEqual(dc2.Labels, dc2.Spec.Template.Labels) {
			if test.expectLabelsChange {
				t.Errorf("expected: %v, got: %v", dc2.Spec.Template.Labels, dc2.Labels)
			} else {
				t.Errorf("unexpected equality: %v", dc.Labels)
			}
		}
	}
}

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	data, err := v1.Codec.Encode(obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := api.Codec.Decode(data)
	if err != nil {
		t.Errorf("%v\nData: %s\nSource: %#v", err, string(data), obj)
		return nil
	}
	obj3 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(runtime.Object)
	err = api.Scheme.Convert(obj2, obj3)
	if err != nil {
		t.Errorf("%v\nSource: %#v", err, obj2)
		return nil
	}
	return obj3
}
