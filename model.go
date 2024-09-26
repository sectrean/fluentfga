package fluentfga

import (
	"cmp"
	"slices"

	proto "github.com/openfga/api/proto/openfga/v1"
)

type Model struct {
	Package string
	Types   []*TypeDefinition
}

type TypeDefinition struct {
	Name string
	Type string

	IDName string
	IDType string

	Relations []*Relation
}

type Relation struct {
	Name string
	Type string

	Object *TypeDefinition

	UserTypes []string
}

func NewModel(
	model *proto.AuthorizationModel,
	config *Config,
) *Model {
	typeMap := make(map[string]*TypeDefinition, len(model.TypeDefinitions))
	types := make([]*TypeDefinition, 0, len(model.TypeDefinitions))

	for _, typeDef := range model.TypeDefinitions {
		td := newTypeDefinition(typeDef, config)
		typeMap[td.Name] = td
		types = append(types, td)
	}

	createRelations(model, config, typeMap)

	return &Model{
		Package: config.Package,
		Types:   types,
	}
}

func newTypeDefinition(
	typeDef *proto.TypeDefinition,
	config *Config,
) *TypeDefinition {
	typ := titleCase(typeDef.Type)

	idName := typ + "ID"
	idType := "string"

	if typeConfig, ok := config.Types[typeDef.Type]; ok {
		typ = typeConfig.Type
		idName = typeConfig.IDName
		idType = typeConfig.IDType
	}

	return &TypeDefinition{
		Name: typeDef.Type,
		Type: typ,

		IDName: idName,
		IDType: idType,
	}
}

func createRelations(
	model *proto.AuthorizationModel,
	_ *Config,
	typeMap map[string]*TypeDefinition,
) {
	for _, typeDef := range model.TypeDefinitions {
		meta := typeDef.Metadata
		td := typeMap[typeDef.Type]

		for name := range typeDef.Relations {
			relMeta := meta.Relations[name]

			userTypes := make([]string, 0, len(relMeta.DirectlyRelatedUserTypes))
			for _, userType := range relMeta.DirectlyRelatedUserTypes {
				ut := titleCase(userType.Type)
				if rel := userType.GetRelation(); rel != "" {
					ut += titleCase(rel) + "Userset"
				}

				userTypes = append(userTypes, ut)
			}

			rel := &Relation{
				Name:      name,
				Type:      titleCase(name),
				Object:    td,
				UserTypes: userTypes,
			}

			td.Relations = append(td.Relations, rel)
		}

		slices.SortFunc(td.Relations, func(a, b *Relation) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}
}
