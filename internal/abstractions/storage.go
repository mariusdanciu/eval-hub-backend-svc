package abstractions

import "github.ibm.com/julpayne/eval-hub-backend-svc/pkg/api"

type Query map[string]string

type Storage interface {
	CreateEvaluationJob(evaluation *api.EvaluationJobResource) error
	GetEvaluationJob(id string) (*api.EvaluationJobResource, error)
	GetEvaluationJobs(query Query) (*api.EvaluationJobResourceList, error)
	DeleteEvaluationJob(id string) error

	CreateCollection(collection *api.CollectionResource) error
	GetCollection(id string) (*api.CollectionResource, error)
	GetCollections(query Query) (*api.CollectionResourceList, error)
	UpdateCollection(collection *api.CollectionResource) error
	DeleteCollection(id string) error
}
