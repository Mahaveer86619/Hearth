package ingestion_pipeline

type IngestionPipelineService struct {
	ingestChan chan []byte
	quit       chan struct{}
}

var instance *IngestionPipelineService

func NewIngestionPipelineService() *IngestionPipelineService {
	if instance == nil {
		instance = &IngestionPipelineService{
			ingestChan: make(chan []byte, 10000),
			quit:       make(chan struct{}),
		}
		return instance
	}
	return instance
}

func GetInstance() *IngestionPipelineService {
	return instance
}
