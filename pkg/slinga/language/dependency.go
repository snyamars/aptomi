package language

import (
	. "github.com/Aptomi/aptomi/pkg/slinga/log"
	log "github.com/Sirupsen/logrus"
)

/*
	This file declares all the necessary structures for Dependencies (User "wants" Service)
*/

// Dependency in a form <UserID> requested <Service> (and provided additional <Labels>)
type Dependency struct {
	*SlingaObject

	Enabled bool
	UserID  string
	Service string
	Labels  map[string]string

	// This fields are populated when dependency gets resolved
	Resolved   bool
	ServiceKey string
}

// UnmarshalYAML is a custom unmarshaller for Dependency, which sets Enabled to True before unmarshalling
func (dependency *Dependency) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Alias Dependency
	instance := Alias{Enabled: true}
	if err := unmarshal(&instance); err != nil {
		return err
	}
	*dependency = Dependency(instance)
	return nil
}

// GlobalDependencies represents the list of global dependencies (see the definition above)
// TODO: during serialization there is data duplication (as both fields get serialized). should prob avoid this
type GlobalDependencies struct {
	// dependencies <service> -> list of dependencies
	DependenciesByService map[string][]*Dependency

	// dependencies <id> -> dependency
	DependenciesByID map[string]*Dependency
}

// NewGlobalDependencies creates and initializes a new empty list of global dependencies
func NewGlobalDependencies() *GlobalDependencies {
	return &GlobalDependencies{
		DependenciesByService: make(map[string][]*Dependency),
		DependenciesByID:      make(map[string]*Dependency),
	}
}

// GetLabelSet applies set of transformations to labels
func (dependency *Dependency) GetLabelSet() LabelSet {
	return LabelSet{Labels: dependency.Labels}
}

// AddDependency appends a single dependency to an existing object
func (src GlobalDependencies) AddDependency(dependency *Dependency) {
	if len(dependency.GetID()) <= 0 {
		Debug.WithFields(log.Fields{
			"dependency": dependency,
		}).Panic("Empty dependency ID")
	}
	src.DependenciesByService[dependency.Service] = append(src.DependenciesByService[dependency.Service], dependency)
	src.DependenciesByID[dependency.GetID()] = dependency
}

// TODO: added temporary method to deal with existing dependency IDs. Once we implement namespaces, may be this has to be re-thinked
func (dependency *Dependency) GetID() string {
	return dependency.GetName()
}

func (dependency *Dependency) GetObjectType() SlingaObjectType {
	return TypePolicy
}
