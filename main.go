package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	populatorMachinery "github.com/kubernetes-csi/lib-volume-populator/populator-machinery"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	prefix    = "k8s.viveksingh.dev"
	mountPath = "/mnt/vol"
)

func main() {
	var image string
	var namespace string
	var mode string
	var uri string

	flag.StringVar(&image, "image", "", "Image for populator component")
	flag.StringVar(&namespace, "namespace", "", "Namespace for populator component")
	flag.StringVar(&mode, "mode", "", "Mode to run the application in")
	flag.StringVar(&uri, "uri", "", "URI for the content of the volume")

	flag.Parse()

	switch mode {
	case "controller":
		gk := schema.GroupKind{
			Group: "k8s.viveksingh.dev",
			Kind:  "GenericHTTPPopulator",
		}
		gvr := schema.GroupVersionResource{
			Group:    gk.Group,
			Version:  "v1alpha1",
			Resource: "generichttppopulators",
		}
		populatorMachinery.RunController("", "", image, "", "", namespace, prefix, gk, gvr, mountPath, "", populatorArgs)
	case "populator":
		// the code that we write here is going to get called from populator pod
		// we know that PVC is mounted at `mountPath`
		populate(uri)
	default:
		log.Printf("Mode %s is not supported", mode)
	}
}

func populate(uri string) {
	if uri == "" {
		log.Printf("URI cannot be empty")
		panic(errors.New("URI can not be empty"))
	}

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Printf("Failed creating new request, error %s\n", err.Error())
		panic(errors.New(err.Error()))
	}

	fileName := path.Base(req.URL.Path)

	f, err := os.Create(path.Join(mountPath, fileName))
	if err != nil {
		log.Printf("failed to create file %s\n", err.Error())
		panic(errors.New(err.Error()))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("failed making the request %s\n", err.Error())
		panic(errors.New(err.Error()))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("failed to read data from response %s\n", err.Error())
		panic(errors.New(err.Error()))
	}

	_, err = f.Write(data)
	if err != nil {
		log.Printf("Failed writing data to file %s\n", err.Error())
		panic(errors.New(err.Error()))
	}

	log.Println("Populated successfully")
}

type GenericHTTPPopulator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec GenericHTTPPopulatorSpec `json:"spec"`
}

type GenericHTTPPopulatorSpec struct {
	URI string `json:"uri"`
}

// b here specifies if the volume is in block mode
// u is the populator instance
func populatorArgs(b bool, u *unstructured.Unstructured) ([]string, error) {
	var ghp GenericHTTPPopulator
	var args []string
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), &ghp)
	if err != nil {
		log.Printf("Failed converting unstructured to GHP, error %s\n", err.Error())
		return args, err
	}

	args = append(args, "--mode=populator")
	args = append(args, fmt.Sprintf("--uri=%s", ghp.Spec.URI))
	return args, nil
}
