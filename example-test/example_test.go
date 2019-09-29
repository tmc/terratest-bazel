package terratest_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/random"
	_ "github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
)

func run(args ...string) {
	a0 := args[0]
	al := args[1:]
	cmd := exec.Command(a0, al...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func kubeconfigPath() string {
	// if TEST_SRCDIR is set we are running under bazel and are sandboxed.
	if tsd := os.Getenv("TEST_SRCDIR"); tsd != "" {
		return tsd + "/kubeconfig/kubeconfig.yaml"
	}
	return ""
}

func TestFoobar(t *testing.T) {
	run("env")
	run("tree")
	run("kubectl", "config", "get-contexts")

	options := k8s.NewKubectlOptions("", kubeconfigPath())
	kubeResourcePath, err := filepath.Abs("./testdata/nginx-deployment.yaml")
	if err != nil {
		t.Error(err)
	}

	namespaceName := fmt.Sprintf("kubernetes-basic-example-%s", strings.ToLower(random.UniqueId()))
	k8s.CreateNamespace(t, options, namespaceName)
	options.Namespace = namespaceName
	// ... and make sure to delete the namespace at the end of the test
	defer k8s.DeleteNamespace(t, options, namespaceName)

	// At the end of the test, run `kubectl delete -f RESOURCE_CONFIG` to clean up any resources that were created.
	defer k8s.KubectlDelete(t, options, kubeResourcePath)

	// This will run `kubectl apply -f RESOURCE_CONFIG` and fail the test if there are any errors
	k8s.KubectlApply(t, options, kubeResourcePath)

	// This will get the service resource and verify that it exists and was retrieved successfully. This function will
	// fail the test if the there is an error retrieving the service resource from Kubernetes.
	service := k8s.GetService(t, options, "nginx-service")
	require.Equal(t, service.Name, "nginx-service")
	expectedName := "nginx-service"
	if service.Name != expectedName {
		t.Errorf("got '%s' and wanted '%s'", service.Name, expectedName)
	}
}
