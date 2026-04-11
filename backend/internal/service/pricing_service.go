package service

import (
	"encoding/json"
	"log"

	"github.com/SManriqueDev/cubearchitect/internal/cubepath"
)

type PricingService struct {
}

func NewPricingService() *PricingService {
	return &PricingService{}
}

func (s *PricingService) GetPlans(client cubepath.CubePathClient) (cubepath.PlansResponse, error) {
	resp, err := client.Get("/vps/plans")
	var plansResp cubepath.PlansResponse
	if err != nil {
		log.Printf("Error fetching plans: %v", err)
		return cubepath.PlansResponse{}, err
	}
	if err := json.Unmarshal(resp, &plansResp); err != nil {
		log.Printf("Error parsing plans: %v", err)
		return cubepath.PlansResponse{}, err
	}

	return plansResp, nil
}

func (s *PricingService) GetPricing(client cubepath.CubePathClient) (json.RawMessage, error) {
	plansRes, err := client.Get("/vps/plans")
	if err != nil {
		log.Printf("Error fetching plans: %v", err)
		return nil, err
	}

	templatesRes, err := client.Get("/vps/templates")
	if err != nil {
		log.Printf("Error fetching templates: %v", err)
		return nil, err
	}

	var plansData cubepath.PlansResponse
	if err := json.Unmarshal(plansRes, &plansData); err != nil {
		log.Printf("Error parsing plans: %v", err)
		return nil, err
	}

	var templatesData cubepath.TemplatesResponse
	if err := json.Unmarshal(templatesRes, &templatesData); err != nil {
		log.Printf("Error parsing templates: %v", err)
		return nil, err
	}

	transformed := transformPricing(plansData, templatesData)
	return json.Marshal(transformed)
}

func transformPricing(plans cubepath.PlansResponse, templates cubepath.TemplatesResponse) map[string]interface{} {
	templatesList := make([]map[string]string, 0)
	seenTemplate := make(map[string]bool)

	for _, os := range templates.OperatingSystems {
		key := os.TemplateName
		if !seenTemplate[key] {
			seenTemplate[key] = true
			templatesList = append(templatesList, map[string]string{
				"template_name": os.TemplateName,
				"os_name":       os.OSName,
				"version":       os.Version,
			})
		}
	}

	result := map[string]interface{}{
		"vps": map[string]interface{}{
			"locations": plans.Locations,
			"templates": templatesList,
		},
	}

	return result
}
