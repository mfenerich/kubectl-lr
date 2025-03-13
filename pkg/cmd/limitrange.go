package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/yaml"
)

var (
	limitExample = `
    # Create a LimitRange with CPU and memory limits in the specified namespace
    kubectl create limitrange my-limitrange --namespace=my-namespace --max-cpu="1" --min-cpu=100m --default-cpu=500m --default-request-cpu=500m --max-memory=500Mi --min-memory=100Mi

    # Create a LimitRange with only CPU limits
    kubectl create limitrange my-cpu-limit --namespace=my-namespace --max-cpu="2" --min-cpu=500m --default-cpu=1 --default-request-cpu=500m --dry-run=client -o yaml
    `
)

// LimitOptions holds information required to create a LimitRange
type LimitOptions struct {
	configFlags       *genericclioptions.ConfigFlags
	namespace         string
	name              string
	maxCPU            string
	minCPU            string
	defaultCPU        string
	defaultRequestCPU string
	maxMemory         string
	minMemory         string
	dryRun            string // Accepts "client" or "server"
	output            string
	IOStreams         genericclioptions.IOStreams

	// Function to create Kubernetes clientset, can be overridden in tests
	clientsetFunc func(config *rest.Config) (kubernetes.Interface, error)
}

// NewLimitOptions initializes an instance of LimitOptions with default values
//
//go:noinline
func NewLimitOptions(streams genericclioptions.IOStreams) *LimitOptions {
	return &LimitOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
		clientsetFunc: func(config *rest.Config) (kubernetes.Interface, error) {
			return kubernetes.NewForConfig(config)
		},
	}
}

// NewCmdCreateLimitRange creates a cobra command wrapping LimitOptions
func NewCmdCreateLimitRange(streams genericiooptions.IOStreams) *cobra.Command {
	o := NewLimitOptions(streams)

	cmd := &cobra.Command{
		Use:          "limitrange NAME [flags]",
		Short:        "Create a LimitRange resource",
		Example:      limitExample,
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			o.name = args[0]
			if err := o.Complete(c, args); err != nil {
				return fmt.Errorf("completion error: %w", err)
			}
			if err := o.Validate(); err != nil {
				return fmt.Errorf("validation error: %w", err)
			}
			if err := o.Run(); err != nil {
				return fmt.Errorf("execution error: %w", err)
			}
			return nil
		},
	}

	// coverage:ignore-start
	// Add common flags
	o.configFlags.AddFlags(cmd.Flags())

	// Add shorthand -n for the --namespace flag
	if nsFlag := cmd.Flag("namespace"); nsFlag != nil {
		nsFlag.Shorthand = "n"
	}

	// Define custom flags for LimitRange options
	cmd.Flags().StringVar(&o.maxCPU, "max-cpu", "", "Maximum CPU limit for containers")
	cmd.Flags().StringVar(&o.minCPU, "min-cpu", "", "Minimum CPU limit for containers")
	cmd.Flags().StringVar(&o.defaultCPU, "default-cpu", "", "Default CPU limit for containers")
	cmd.Flags().StringVar(&o.defaultRequestCPU, "default-request-cpu", "", "Default CPU request for containers")
	cmd.Flags().StringVar(&o.maxMemory, "max-memory", "", "Maximum memory limit for containers")
	cmd.Flags().StringVar(&o.minMemory, "min-memory", "", "Minimum memory limit for containers")
	cmd.Flags().StringVar(&o.dryRun, "dry-run", "", "Must be 'client' or 'server'. If set, only print the object that would be sent without sending it.")
	cmd.Flags().StringVarP(&o.output, "output", "o", "", "Output format. One of: yaml|json")

	return cmd
	// coverage:ignore-end
}

// Complete sets all required information for creating a LimitRange
func (o *LimitOptions) Complete(_ *cobra.Command, _ []string) error {
	if o.namespace == "" {
		var err error
		o.namespace, _, err = o.configFlags.ToRawKubeConfigLoader().Namespace()
		if err != nil {
			return fmt.Errorf("failed to get current namespace: %w", err)
		}
	}
	return nil
}

// Validate checks that all required arguments and flag values are provided
func (o *LimitOptions) Validate() error {
	if o.namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if o.name == "" {
		return fmt.Errorf("name is required")
	}
	if o.maxMemory == "" && o.minMemory == "" && o.maxCPU == "" && o.minCPU == "" && o.defaultCPU == "" && o.defaultRequestCPU == "" {
		return fmt.Errorf("at least one resource limit or request must be specified")
	}

	resourceFields := map[string]string{
		"max-cpu":             o.maxCPU,
		"min-cpu":             o.minCPU,
		"default-cpu":         o.defaultCPU,
		"default-request-cpu": o.defaultRequestCPU,
		"max-memory":          o.maxMemory,
		"min-memory":          o.minMemory,
	}

	for fieldName, value := range resourceFields {
		if value != "" {
			quantity, err := resource.ParseQuantity(value)
			if err != nil {
				return fmt.Errorf("invalid %s value: %s", fieldName, err)
			}
			if quantity.Sign() != 1 {
				// Sign() returns -1 for negative, 0 for zero, 1 for positive
				return fmt.Errorf("invalid %s value: must be greater than zero", fieldName)
			}
		}
	}
	return nil
}

// Run executes the creation of the LimitRange or prints the YAML/JSON
func (o *LimitOptions) Run() error {
	// Add validation
	if err := o.Validate(); err != nil {
		return err
	}

	limitRange := o.createLimitRangeObject()

	// Handle client-side dry-run
	if o.dryRun == "client" {
		return o.printOutputWithTypeMeta(limitRange)
	} else if o.dryRun != "" && o.dryRun != "server" {
		return fmt.Errorf("invalid value for --dry-run: %s, must be 'client' or 'server'", o.dryRun)
	}

	// Set CreateOptions for server-side dry-run
	createOptions := metav1.CreateOptions{}
	if o.dryRun == "server" {
		createOptions.DryRun = []string{"All"}
	}

	config, err := o.configFlags.ToRawKubeConfigLoader().ClientConfig()
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client config: %w", err)
	}

	clientset, err := o.clientsetFunc(config)
	if err != nil {
		return fmt.Errorf("failed to create Kubernetes clientset: %w", err)
	}

	// Execute the create operation with the given options
	createdLimitRange, err := clientset.CoreV1().LimitRanges(o.namespace).Create(context.TODO(), limitRange, createOptions)
	if err != nil {
		return fmt.Errorf("failed to create LimitRange: %w", err)
	}

	// Print server response for server-side dry-run
	if o.dryRun == "server" {
		if createdLimitRange == nil {
			// Use the original limitRange if createdLimitRange is nil
			createdLimitRange = limitRange
		}
		return o.printOutputWithTypeMeta(createdLimitRange)
	}

	// Print success message
	fmt.Fprintf(o.IOStreams.Out, "limitrange.core %q created\n", limitRange.Name)
	return nil
}

// createLimitRangeObject creates a new LimitRange object populated with provided options
func (o *LimitOptions) createLimitRangeObject() *v1.LimitRange {
	limitRange := &v1.LimitRange{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "LimitRange",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      o.name,
			Namespace: o.namespace,
		},
		Spec: v1.LimitRangeSpec{
			Limits: []v1.LimitRangeItem{
				{
					Type:           v1.LimitTypeContainer,
					Max:            v1.ResourceList{},
					Min:            v1.ResourceList{},
					Default:        v1.ResourceList{},
					DefaultRequest: v1.ResourceList{},
				},
			},
		},
	}

	// Populate resource values if they are provided
	if o.maxCPU != "" {
		limitRange.Spec.Limits[0].Max[v1.ResourceCPU] = resource.MustParse(o.maxCPU)
	}
	if o.minCPU != "" {
		limitRange.Spec.Limits[0].Min[v1.ResourceCPU] = resource.MustParse(o.minCPU)
	}
	if o.defaultCPU != "" {
		limitRange.Spec.Limits[0].Default[v1.ResourceCPU] = resource.MustParse(o.defaultCPU)
	}
	if o.defaultRequestCPU != "" {
		limitRange.Spec.Limits[0].DefaultRequest[v1.ResourceCPU] = resource.MustParse(o.defaultRequestCPU)
	}
	if o.maxMemory != "" {
		limitRange.Spec.Limits[0].Max[v1.ResourceMemory] = resource.MustParse(o.maxMemory)
	}
	if o.minMemory != "" {
		limitRange.Spec.Limits[0].Min[v1.ResourceMemory] = resource.MustParse(o.minMemory)
	}

	return limitRange
}

// printOutputWithTypeMeta ensures TypeMeta is set and prints the LimitRange in the specified format
func (o *LimitOptions) printOutputWithTypeMeta(limitRange *v1.LimitRange) error {
	// Ensure TypeMeta is set
	if limitRange.TypeMeta.APIVersion == "" || limitRange.TypeMeta.Kind == "" {
		limitRange.TypeMeta = metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "LimitRange",
		}
	}

	var output []byte
	var err error

	if o.output == "yaml" {
		output, err = yaml.Marshal(limitRange)
	} else if o.output == "json" {
		serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, nil, nil, json.SerializerOptions{Pretty: true})
		output, err = runtime.Encode(serializer, limitRange)
	} else {
		return fmt.Errorf("unsupported output format: %s", o.output)
	}

	if err != nil {
		return fmt.Errorf("failed to format output: %w", err)
	}

	fmt.Fprintf(o.IOStreams.Out, "%s\n", output)
	return nil
}
