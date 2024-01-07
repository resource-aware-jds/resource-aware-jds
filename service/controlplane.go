package service

import (
	"context"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/resource-aware-jds/resource-aware-jds/config"
	"github.com/resource-aware-jds/resource-aware-jds/models"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/cert"
	"github.com/resource-aware-jds/resource-aware-jds/pkg/distribution"
	"github.com/resource-aware-jds/resource-aware-jds/repository"
	"time"
)

var (
	ErrNodeAlreadyRegistered = errors.New("there is the other node registered with the inputted key")
)

type ControlPlane struct {
	nodeRegistryRepository repository.INodeRegistry
	jobRepository          repository.IJob
	taskRepository         repository.ITask
	caCertificate          cert.CACertificate
	config                 config.ControlPlaneConfigModel
}

type IControlPlane interface {
	RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.KeyData) (certificate cert.TLSCertificate, err error)
	GetAllWorkerNodeFromRegistry(ctx context.Context) ([]models.NodeEntry, error)
	CreateJob(ctx context.Context, imageURL string, taskAttributes [][]byte) (*models.Job, []models.Task, error)
	GetAvailableTask(ctx context.Context) ([]models.Task, error)
	UpdateTaskAfterDistribution(ctx context.Context, successTask []models.Task, errorTask []distribution.DistributeError) error
}

func ProvideControlPlane(jobRepository repository.IJob, taskRepository repository.ITask, nodeRegistryRepository repository.INodeRegistry, caCertificate cert.CACertificate, config config.ControlPlaneConfigModel) IControlPlane {
	return &ControlPlane{
		jobRepository:          jobRepository,
		nodeRegistryRepository: nodeRegistryRepository,
		taskRepository:         taskRepository,
		caCertificate:          caCertificate,
		config:                 config,
	}
}

func (s *ControlPlane) RegisterWorker(ctx context.Context, ip string, port int32, nodePublicKey cert.KeyData) (certificate cert.TLSCertificate, err error) {
	hashedPublicKey, err := nodePublicKey.GetSHA1Hash()
	if err != nil {
		return nil, err
	}

	isExists, err := s.nodeRegistryRepository.IsNodeAlreadyRegistered(ctx, hashedPublicKey)
	if err != nil {
		return nil, err
	}
	if isExists {
		return nil, ErrNodeAlreadyRegistered
	}

	// Sign the certificate.
	clientUUID := uuid.New()
	signedCertificate, err := s.caCertificate.CreateCertificateAndSign(
		pkix.Name{
			CommonName:   fmt.Sprintf("RAJDS Worker %s", clientUUID.String()),
			SerialNumber: clientUUID.String(),
		},
		nodePublicKey,
		365*24*time.Hour,
	)
	if err != nil {
		return nil, err
	}

	// Insert the certificate in the database.
	err = s.nodeRegistryRepository.RegisterWorkerNodeWithCertificate(ctx, ip, port, signedCertificate)
	if err != nil {
		return nil, err
	}

	// Response the certificate back.
	return signedCertificate, nil
}

func (s *ControlPlane) CreateJob(ctx context.Context, imageURL string, taskAttributes [][]byte) (*models.Job, []models.Task, error) {
	// Create Job
	job := models.Job{
		Status:   models.PendingJobStatus,
		ImageURL: imageURL,
	}
	insertedJobID, err := s.jobRepository.Insert(ctx, job)
	if err != nil {
		return nil, nil, err
	}
	job.ID = insertedJobID

	// Create Tasks
	tasks := make([]models.Task, 0, len(taskAttributes))
	for _, taskAttribute := range taskAttributes {
		newTask := models.Task{
			JobID:          insertedJobID,
			Status:         models.CreatedTaskStatus,
			ImageUrl:       imageURL,
			TaskAttributes: taskAttribute,
		}
		tasks = append(tasks, newTask)
	}
	err = s.taskRepository.InsertMany(ctx, tasks)
	if err != nil {
		return nil, nil, err
	}

	tasksResponse, err := s.taskRepository.FindManyByJobID(ctx, insertedJobID)
	if err != nil {
		return nil, nil, err
	}

	return &job, tasksResponse, nil
}

func (s *ControlPlane) GetAllWorkerNodeFromRegistry(ctx context.Context) ([]models.NodeEntry, error) {
	return s.nodeRegistryRepository.GetAllWorkerNode(ctx)
}

func (s *ControlPlane) GetAvailableTask(ctx context.Context) ([]models.Task, error) {
	return s.taskRepository.GetTaskToDistribute(ctx)
}

func (s *ControlPlane) UpdateTaskAfterDistribution(ctx context.Context, successTasks []models.Task, errorTasks []distribution.DistributeError) error {
	taskToUpdate := make([]models.Task, 0, len(successTasks)+len(errorTasks))
	taskToUpdate = append(taskToUpdate, successTasks...)

	// Add failure task
	for _, errorTask := range errorTasks {
		taskToUpdate = append(taskToUpdate, errorTask.Task)
	}

	return s.taskRepository.BulkWriteStatusAndLogByID(ctx, taskToUpdate)
}
