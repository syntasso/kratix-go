package main

import (
	"log"

	kratix "github.com/syntasso/kratix-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	sdk := kratix.New()

	log.Println("Helper variables:")
	log.Printf("Workflow action: %s\n", sdk.WorkflowAction())
	log.Printf("Workflow type: %s\n", sdk.WorkflowType())
	log.Printf("Promise name: %s\n", sdk.PromiseName())
	log.Printf("Pipeline name: %s\n", sdk.PipelineName())

	log.Println("\nReading resource input...")
	resource, err := sdk.ReadResourceInput()
	if err != nil {
		log.Fatalf("failed to read resource input: %v", err)
	}

	log.Println("Getting fields...")
	fields, err := resource.GetValue("spec.fields")
	if err != nil {
		log.Fatalf("failed to get fields: %v", err)
	}
	log.Printf("Fields: %v\n", fields)

	log.Println("Creating configmap...")
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName() + "-config",
			Namespace: resource.GetNamespace(),
		},
		Data: make(map[string]string),
	}

	fieldsMap := fields.([]any)
	for _, field := range fieldsMap {
		fieldMap := field.(map[string]any)
		cm.Data[fieldMap["name"].(string)] = fieldMap["value"].(string)
	}
	cm.Data["workflowAction"] = sdk.WorkflowAction()
	cm.Data["workflowType"] = sdk.WorkflowType()
	cm.Data["promiseName"] = sdk.PromiseName()
	cm.Data["pipelineName"] = sdk.PipelineName()

	cmContent, err := yaml.Marshal(cm)
	if err != nil {
		log.Fatalf("failed to marshal config map: %v", err)
	}
	log.Printf("Config map content: %s\n", string(cmContent))

	log.Println("Writing output...")
	err = sdk.WriteOutput("config.yaml", cmContent)
	if err != nil {
		log.Fatalf("failed to write output: %v", err)
	}

	log.Println("Publishing 'publishedDirectly' status...")
	publishedStatus := kratix.NewStatus()
	publishedStatus.Set("publishedDirectly", true)
	err = sdk.PublishStatus(resource, publishedStatus)
	if err != nil {
		log.Fatalf("failed to publish status: %v", err)
	}

	log.Println("Writing status file...")
	status := kratix.NewStatus()
	status.Set("viaFile", true)
	err = sdk.WriteStatus(status)
	if err != nil {
		log.Fatalf("failed to write status: %v", err)
	}

	log.Println("Validating the Status file...")
	statusFromFile, err := sdk.ReadStatus()
	if err != nil {
		log.Fatalf("failed to read status: %v", err)
	}
	log.Printf("Status from file: %v\n", statusFromFile)

	if statusFromFile.Get("viaFile") != true {
		log.Fatalf("either writing or reading the status file failed")
	}

	log.Println("writing destination selectors file...")
	destinationSelectors := []kratix.DestinationSelector{
		{MatchLabels: map[string]string{"environment": "dev"}},
	}
	err = sdk.WriteDestinationSelectors(destinationSelectors)
	if err != nil {
		log.Fatalf("failed to write destination selectors: %v", err)
	}
	log.Println("Validating the Destination Selectors file...")
	destinationSelectorsFromFile, err := sdk.ReadDestinationSelectors()
	if err != nil {
		log.Fatalf("failed to read destination selectors: %v", err)
	}
	log.Printf("Destination selectors from file: %v\n", destinationSelectorsFromFile)

	log.Println("Validating the Destination Selectors file...")
	if len(destinationSelectorsFromFile) != 1 || destinationSelectorsFromFile[0].MatchLabels["environment"] != "dev" {
		log.Fatalf("either writing or reading the destination selectors file failed")
	}

	log.Println("All tests passed")
}
