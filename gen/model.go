package gen

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

	Usersets    []*Userset
	HasWildcard bool
}

type Relation struct {
	Name string
	Type string

	Object *TypeDefinition

	UserTypes []string
}

type Userset struct {
	Type     string
	Relation string

	Object *TypeDefinition
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

	m := &Model{
		Package: config.Package,
		Types:   types,
	}
	m.sort()

	return m
}

func newTypeDefinition(
	typeDef *proto.TypeDefinition,
	config *Config,
) *TypeDefinition {
	name := typeDef.Type
	td := &TypeDefinition{
		Name:   name,
		Type:   config.FormatTypeName(name),
		IDName: config.FormatIDName(name),
		IDType: "string",
	}

	if tc, ok := config.Types[typeDef.Type]; ok {
		if tc.Type != "" {
			td.Type = tc.Type
		}
		if tc.IDName != "" {
			td.IDName = tc.IDName
		}
		if tc.IDType != "" {
			td.IDType = tc.IDType
		}
	}

	return td
}

func createRelations(
	model *proto.AuthorizationModel,
	config *Config,
	typeMap map[string]*TypeDefinition,
) {
	for _, typeDef := range model.TypeDefinitions {
		meta := typeDef.Metadata
		td := typeMap[typeDef.Type]

		for name := range typeDef.Relations {
			relMeta := meta.Relations[name]

			userTypes := make([]string, 0, len(relMeta.DirectlyRelatedUserTypes))
			for _, userType := range relMeta.DirectlyRelatedUserTypes {
				ut := config.TypeName(userType.Type)
				usType := typeMap[userType.Type]

				if rel := userType.GetRelation(); rel != "" {
					ut += config.FormatTypeName(rel) + "Userset"
					us := &Userset{
						Type:     ut,
						Relation: rel,
						Object:   usType,
					}
					usType.Usersets = append(usType.Usersets, us)
				} else if rel := userType.GetWildcard(); rel != nil {
					ut += config.FormatTypeName(rel.String()) + "Wildcard"

					usType.HasWildcard = true
				}

				userTypes = append(userTypes, ut)
			}

			rel := &Relation{
				Name:      name,
				Type:      config.TypeName(name),
				Object:    td,
				UserTypes: userTypes,
			}

			td.Relations = append(td.Relations, rel)
		}
	}
}

func (m *Model) sort() {
	for _, td := range m.Types {
		// Sort and de-dupe usersets
		slices.SortFunc(td.Usersets, func(a, b *Userset) int {
			return cmp.Or(
				cmp.Compare(a.Type, b.Type),
				cmp.Compare(a.Relation, b.Relation),
			)
		})

		td.Usersets = slices.CompactFunc(td.Usersets, func(a, b *Userset) bool {
			return a.Type == b.Type && a.Relation == b.Relation
		})

		// Sort relations
		slices.SortFunc(td.Relations, func(a, b *Relation) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}
}
