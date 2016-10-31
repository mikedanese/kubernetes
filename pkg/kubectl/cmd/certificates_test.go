/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"bytes"
	"net/http"
	"testing"

	"k8s.io/kubernetes/pkg/apis/certificates"
	"k8s.io/kubernetes/pkg/client/unversioned/fake"
	cmdtesting "k8s.io/kubernetes/pkg/kubectl/cmd/testing"
)

func TestCertificateApprove(t *testing.T) {
	csr := &certificates.CertificateSigningRequest{}
	f, tf, codec, ns := cmdtesting.NewAPIFactory()
	tf.Printer = &testPrinter{}
	tf.Client = &fake.RESTClient{
		NegotiatedSerializer: ns,
		Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			switch p, m := req.URL.Path, req.Method; {
			case p == "/certificatesigningrequests/foo" && m == "GET":
				return &http.Response{StatusCode: 200, Header: defaultHeader(), Body: objBody(codec, csr)}, nil
			case p == "/certificatesigningrequests/foo/approval" && m == "POST":
				return &http.Response{StatusCode: 201, Header: defaultHeader(), Body: objBody(codec, csr)}, nil
			default:
				t.Fatalf("unexpected request: %#v\n%#v", req.URL, req)
				return nil, nil
			}
		}),
	}
	//tf.ClientConfig = &restclient.Config{
	//	APIPath: "/api",
	//	ContentConfig: restclient.ContentConfig{
	//		NegotiatedSerializer: api.Codecs,
	//		GroupVersion:         &registered.GroupOrDie(api.GroupName).GroupVersion,
	//	},
	//}
	tf.Namespace = "test"
	buf := bytes.NewBuffer([]byte{})
	cmd := NewCmdCertificateApprove(f, buf)
	//cmd.Flags().Set("output", "name")
	cmd.Run(cmd, []string{"foo"})
	//expectedOutput := "serviceaccount/" + serviceAccountObject.Name + "\n"
	//if buf.String() != expectedOutput {
	//	t.Errorf("expected output: %s, but got: %s", expectedOutput, buf.String())
	//}
}
