package update

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
	"log"
	"os"
	"strings"
)

var (
	project         = ""
	watchedRegistry = ""
	endpoint        = ""
)

type PubSubMessage struct {
	Data string `json:"data"`
}

type EventData struct {
	Action string `json:"action"`
	Digest string `json:"digest"`
	Tag    string `json:"tag"`
}

func init() {
	project = os.Getenv("PROJECT")
	watchedRegistry = os.Getenv("REGISTRY")
	endpoint = os.Getenv("ENDPOINT")
}

func Update(ctx context.Context, m PubSubMessage) error {
	raw, _ := base64.StdEncoding.DecodeString(m.Data)
	var event EventData
	err := json.Unmarshal(raw, &event)

	if err != nil {
		return fmt.Errorf("could not extract event: %w", err)
	}

	if event.Action != "INSERT" {
		log.Printf("action is not \"INSERT\". Skipping event")
		return nil
	}

	if !strings.HasPrefix(event.Tag, watchedRegistry) {
		log.Printf("tag \"%s\" is not from watched registry \"%s\". Skipping event\n", event.Tag, watchedRegistry)
		return nil
	}

	serviceName := extractServiceName(event.Tag)

	svc, err := run.NewService(ctx, option.WithEndpoint(endpoint))
	if err != nil {
		return fmt.Errorf("could not create api client: %w", err)
	}

	name := fmt.Sprintf("namespaces/%s/services/%s", project, serviceName)

	service, err := svc.Namespaces.Services.Get(name).Do()
	if err != nil {
		return fmt.Errorf("could not get project: %w", err)
	}

	service.Spec.Template.Metadata.Name = fmt.Sprintf("%s-%05d", serviceName, service.Metadata.Generation)
	service.Spec.Template.Spec.Containers[0].Image = event.Digest

	_, err = svc.Namespaces.Services.ReplaceService(name, service).Do()
	if err != nil {
		return fmt.Errorf("could not update service: %w", err)
	}

	return nil
}

// extractServiceName takes the newly created tag and extracts the service name from it.
func extractServiceName(tag string) string {
	serviceName := strings.Replace(tag, watchedRegistry, "", 1)
	serviceName = strings.Replace(serviceName, ":", "-", 1)
	serviceName = strings.Replace(serviceName, "/", "", -1)
	return serviceName
}
