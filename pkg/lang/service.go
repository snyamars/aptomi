package lang

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/object"
	"github.com/Aptomi/aptomi/pkg/util"
	"sync"
)

// ServiceObject is an informational data structure with Kind and Constructor for Service
var ServiceObject = &object.Info{
	Kind:        "service",
	Versioned:   true,
	Constructor: func() object.Base { return &Service{} },
}

// Service defines individual service
type Service struct {
	Metadata

	Labels     map[string]string
	Components []*ServiceComponent

	// Lazily evaluated fields (all components topologically sorted). Use via getter
	componentsOrderedOnce sync.Once
	componentsOrderedErr  error
	componentsOrdered     []*ServiceComponent

	componentsMapOnce sync.Once
	componentsMap     map[string]*ServiceComponent
}

// ServiceComponent defines component within a service
type ServiceComponent struct {
	Name string `validate:"identifier"`

	// Contract, if not empty, means that component points to another contract as a dependency
	Contract string `validate:"identifier"`

	// Code, if not empty, means that component is a code that can be instantiated
	Code         *Code
	Discovery    util.NestedParameterMap
	Dependencies []string `validate:"identifier"`
}

// Code with type and parameters, used to instantiate/update/delete component instances
type Code struct {
	Type   string `validate:"codetype"`
	Params util.NestedParameterMap
}

// GetComponentsMap lazily initializes and returns a map of name -> component
// This should be thread safe
func (service *Service) GetComponentsMap() map[string]*ServiceComponent {
	service.componentsMapOnce.Do(func() {
		// Put all components into map
		service.componentsMap = make(map[string]*ServiceComponent)
		for _, c := range service.Components {
			service.componentsMap[c.Name] = c
		}
	})
	return service.componentsMap
}

// Topologically sort components of a given service and return error if there is a cycle detected
func (service *Service) dfsComponentSort(u *ServiceComponent, colors map[string]int) error {
	colors[u.Name] = 1

	for _, vName := range u.Dependencies {
		v, exists := service.GetComponentsMap()[vName]
		if !exists {
			return fmt.Errorf("service '%s' has a dependency on non-existing component '%s'", service.Name, vName)
		}
		if vColor, ok := colors[v.Name]; !ok {
			// not visited yet -> visit and exit if a cycle was found or another error occurred
			if err := service.dfsComponentSort(v, colors); err != nil {
				return err
			}
		} else if vColor == 1 {
			return fmt.Errorf("component cycle detected while processing service '%s' component '%s'", service.Name, vName)
		}
	}

	service.componentsOrdered = append(service.componentsOrdered, u)
	colors[u.Name] = 2
	return nil
}

// GetComponentsSortedTopologically returns all components sorted in a topological order
// This should be thread safe
func (service *Service) GetComponentsSortedTopologically() ([]*ServiceComponent, error) {
	service.componentsOrderedOnce.Do(func() {
		// Initiate colors
		colors := make(map[string]int)

		// Dfs
		for _, c := range service.Components {
			if _, ok := colors[c.Name]; !ok {
				if err := service.dfsComponentSort(c, colors); err != nil {
					service.componentsOrdered = nil
					service.componentsOrderedErr = err
					break
				}
			}
		}
	})

	return service.componentsOrdered, service.componentsOrderedErr
}
