package main

import (
	"log"

	"github.com/syntasso/kratix-go"
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

	log.Println("\nReading promise input...")
	promise, err := sdk.ReadPromiseInput()
	if err != nil {
		log.Fatalf("failed to read promise input: %v", err)
	}

	log.Println("Creating configmap...")
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      promise.GetName() + "-config",
			Namespace: "default",
		},
		Data: make(map[string]string),
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

	log.Println("All tests passed")
}
