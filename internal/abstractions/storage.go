package abstractions

import "github.ibm.com/julpayne/eval-hub-backend-svc/pkg/api"

type Query map[string]string

type Storage interface {
	CreateEvaluationJOb(evaluation *api.EvaluationJobResource) error
	GetEvaluation(id string) (*api.EvaluationJobResource, error)
	GetEvaluations(query Query) (*api.EvaluationJobResourceList, error)
	UpdateEvaluation(evaluation *api.EvaluationJobResource) error
	DeleteEvaluation(id string) error

	CreateCollection(collection *api.CollectionResource) error
	GetCollection(id string) (*api.CollectionResource, error)
	GetCollections(query Query) (*api.CollectionResourceList, error)
	UpdateCollection(collection *api.CollectionResource) error
	DeleteCollection(id string) error
}
