package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
)

func TestCreateLimitRangeObject(t *testing.T) {
    options := &LimitOptions{
        name:      "test-limitrange",
        namespace: "default",
        maxCPU:    "1",
        minCPU:    "100m",
        maxMemory: "512Mi",
        minMemory: "128Mi",
    }

    limitRange := options.createLimitRangeObject()

    assert.Equal(t, "test-limitrange", limitRange.ObjectMeta.Name)
    assert.Equal(t, "default", limitRange.ObjectMeta.Namespace)

    // Create temporary variables for map values
    maxCPU := limitRange.Spec.Limits[0].Max[v1.ResourceCPU]
    minCPU := limitRange.Spec.Limits[0].Min[v1.ResourceCPU]
    maxMemory := limitRange.Spec.Limits[0].Max[v1.ResourceMemory]
    minMemory := limitRange.Spec.Limits[0].Min[v1.ResourceMemory]

    // Call String() on pointers to the quantities
    assert.Equal(t, "1", (&maxCPU).String())
    assert.Equal(t, "100m", (&minCPU).String())
    assert.Equal(t, "512Mi", (&maxMemory).String())
    assert.Equal(t, "128Mi", (&minMemory).String())
}

func TestComplete(t *testing.T) {
    options := &LimitOptions{
        configFlags: genericclioptions.NewConfigFlags(true),
    }

    cmd := &cobra.Command{}
    args := []string{"test-limitrange"}

    // Set a namespace in configFlags
    options.configFlags.Namespace = new(string)
    *options.configFlags.Namespace = "test-namespace"

    err := options.Complete(cmd, args)
    assert.NoError(t, err, "expected no error during completion")
    assert.Equal(t, "test-namespace", options.namespace)
}

func TestValidate(t *testing.T) {
    // Initialize with at least one resource limit
    options := &LimitOptions{
        namespace: "default",
        name:      "test-limitrange",
        maxCPU:    "1", // Added resource limit
    }

    // Should pass validation
    err := options.Validate()
    assert.NoError(t, err, "expected no error for valid input")

    // Test for empty namespace
    options.namespace = ""
    err = options.Validate()
    assert.Error(t, err, "expected error when namespace is empty")
    options.namespace = "default" // Reset namespace

    // Test for missing name
    options.name = ""
    err = options.Validate()
    assert.Error(t, err, "expected error when name is missing")
}

func TestValidateInvalidResourceQuantities(t *testing.T) {
    options := &LimitOptions{
        namespace: "default",
        name:      "test-limitrange",
        maxCPU:    "invalid", // Invalid resource quantity
    }

    err := options.Validate()
    assert.Error(t, err, "expected error for invalid maxCPU value")
    assert.Contains(t, err.Error(), "invalid max-cpu value")
}

func TestValidateInvalidResourceFormat(t *testing.T) {
    options := &LimitOptions{
        namespace: "default",
        name:      "test-limitrange",
        maxCPU:    "-1", // Negative value
    }

    err := options.Validate()
    if assert.Error(t, err, "expected error for negative maxCPU value") {
        assert.Contains(t, err.Error(), "invalid max-cpu value: must not be negative")
    }
}

func TestPrintOutputWithTypeMeta(t *testing.T) {
    options := &LimitOptions{
        output: "yaml",
        IOStreams: genericclioptions.IOStreams{
            Out: new(bytes.Buffer),
        },
    }

    limitRange := options.createLimitRangeObject()
    err := options.printOutputWithTypeMeta(limitRange)
    assert.NoError(t, err, "expected no error while printing output")

    // Verify output
    output := options.IOStreams.Out.(*bytes.Buffer).String()
    assert.Contains(t, output, "apiVersion: v1")
    assert.Contains(t, output, "kind: LimitRange")
}

func TestPrintOutputWithTypeMetaJSON(t *testing.T) {
    options := &LimitOptions{
        output: "json",
        IOStreams: genericclioptions.IOStreams{
            Out: new(bytes.Buffer),
        },
    }

    limitRange := options.createLimitRangeObject()
    err := options.printOutputWithTypeMeta(limitRange)
    assert.NoError(t, err, "expected no error while printing output")

    // Verify output
    output := options.IOStreams.Out.(*bytes.Buffer).String()
    assert.Contains(t, output, "\"apiVersion\": \"v1\"")
    assert.Contains(t, output, "\"kind\": \"LimitRange\"")
}

func TestRunWithFakeClient(t *testing.T) {
    fakeClientset := fake.NewSimpleClientset()

    options := &LimitOptions{
        name:        "test-limitrange",
        namespace:   "default",
        maxCPU:      "1",
        IOStreams:   genericclioptions.IOStreams{Out: new(bytes.Buffer)},
        configFlags: genericclioptions.NewConfigFlags(true),
        clientsetFunc: func(config *rest.Config) (kubernetes.Interface, error) {
            return fakeClientset, nil
        },
    }

    // Mock ToRESTConfig to return a dummy config
    options.configFlags.WrapConfigFn = func(c *rest.Config) *rest.Config {
        return &rest.Config{}
    }

    // Run the method
    err := options.Run()
    assert.NoError(t, err, "expected no error during Run")

    // Verify that the LimitRange was created
    lr, err := fakeClientset.CoreV1().LimitRanges(options.namespace).Get(context.TODO(), options.name, metav1.GetOptions{})
    assert.NoError(t, err, "expected no error getting LimitRange")
    assert.Equal(t, options.name, lr.Name)
}

func TestRunDryRunClient(t *testing.T) {
    options := &LimitOptions{
        name:      "test-limitrange",
        namespace: "default",
        maxCPU:    "1",
        dryRun:    "client",
        output:    "yaml",
        IOStreams: genericclioptions.IOStreams{
            Out: new(bytes.Buffer),
        },
    }

    err := options.Run()
    assert.NoError(t, err, "expected no error during Run with dry-run=client")

    // Verify output
    output := options.IOStreams.Out.(*bytes.Buffer).String()
    assert.Contains(t, output, "apiVersion: v1")
    assert.Contains(t, output, "kind: LimitRange")
    assert.Contains(t, output, "name: test-limitrange")
}

func TestRunDryRunServer(t *testing.T) {
    fakeClientset := fake.NewSimpleClientset()

    // Flag to indicate dry-run is enabled
    dryRunEnabled := false

    options := &LimitOptions{
        name:        "test-limitrange",
        namespace:   "default",
        maxCPU:      "1",
        dryRun:      "server",
        output:      "yaml", // Set output format
        IOStreams:   genericclioptions.IOStreams{Out: new(bytes.Buffer)},
        configFlags: genericclioptions.NewConfigFlags(true),
        clientsetFunc: func(config *rest.Config) (kubernetes.Interface, error) {
            return fakeClientset, nil
        },
    }

    // Set dryRunEnabled based on options
    dryRunEnabled = options.dryRun == "server"

    // Add reactor to simulate DryRun
    fakeClientset.Fake.PrependReactor("create", "limitranges", func(action k8stesting.Action) (bool, runtime.Object, error) {
        createAction := action.(k8stesting.CreateAction)
        limitRange := createAction.GetObject().(*v1.LimitRange)

        if dryRunEnabled {
            // Simulate DryRun by not storing the object
            return true, limitRange, nil
        }
        // Proceed with normal behavior
        return false, nil, nil
    })

    // Mock ToRESTConfig to return a dummy config
    options.configFlags.WrapConfigFn = func(c *rest.Config) *rest.Config {
        return &rest.Config{}
    }

    err := options.Run()
    assert.NoError(t, err, "expected no error during Run with dry-run=server")

    // Verify that the LimitRange was not actually created
    _, err = fakeClientset.CoreV1().LimitRanges(options.namespace).Get(context.TODO(), options.name, metav1.GetOptions{})
    assert.Error(t, err, "expected error getting LimitRange, since it should not be created in dry-run=server mode")
}

func TestRunWithInvalidDryRunOption(t *testing.T) {
    options := &LimitOptions{
        name:      "test-limitrange",
        namespace: "default",
        maxCPU:    "1",
        dryRun:    "invalid",
        IOStreams: genericclioptions.IOStreams{
            Out: new(bytes.Buffer),
        },
    }

    err := options.Run()
    assert.Error(t, err, "expected error due to invalid dry-run option")
    assert.Contains(t, err.Error(), "invalid value for --dry-run")
}

func TestRunWithOutputOption(t *testing.T) {
    options := &LimitOptions{
        name:      "test-limitrange",
        namespace: "default",
        maxCPU:    "1",
        output:    "yaml",
        dryRun:    "client",
        IOStreams: genericclioptions.IOStreams{
            Out: new(bytes.Buffer),
        },
    }

    err := options.Run()
    assert.NoError(t, err, "expected no error during Run with output option")

    // Verify output
    output := options.IOStreams.Out.(*bytes.Buffer).String()
    assert.Contains(t, output, "apiVersion: v1")
    assert.Contains(t, output, "kind: LimitRange")
    assert.Contains(t, output, "name: test-limitrange")
}

func TestValidateNegativeMinCPU(t *testing.T) {
    options := &LimitOptions{
        namespace: "default",
        name:      "test-limitrange",
        minCPU:    "-100m",
    }

    err := options.Validate()
    if assert.Error(t, err, "expected error for negative minCPU value") {
        assert.Contains(t, err.Error(), "invalid min-cpu value: must not be negative")
    }
}

