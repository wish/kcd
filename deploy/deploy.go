package deploy

import (
	"fmt"
	"log"
	"strings"
	"time"

	conf "github.com/nearmap/cvmanager/config"
	"github.com/nearmap/cvmanager/events"
	cv1 "github.com/nearmap/cvmanager/gok8s/apis/custom/v1"
	"github.com/nearmap/cvmanager/history"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// RolloutTarget defines an interface for something deployable, such as a Deployment, DaemonSet, Pod, etc.
type RolloutTarget interface {
	// Name is the name of the workload (without the namespace).
	Name() string

	// Namespace returns the namespace the workload belongs to.
	Namespace() string

	// Type returns the type of the spec.
	Type() string

	// PodSpec returns the PodSpec for the workload.
	PodSpec() corev1.PodSpec

	// PatchPodSpec receives a pod spec and container which is to be patched
	// according to an appropriate strategy for the type.
	PatchPodSpec(cv *cv1.ContainerVersion, container corev1.Container, version string) error

	// RollbackAfter indicates duration after which a failed rollout
	// should attempt rollback
	RollbackAfter() *time.Duration

	//ProgressHealth indicates weather the current status of progress healthy or not
	ProgressHealth() bool
}

// TemplateRolloutTarget defines methods for deployable resources that manage a collection
// of pods via a pod template. More deployment options are available for such
// resources.
type TemplateRolloutTarget interface {
	RolloutTarget

	// PodTemplateSpec returns the PodTemplateSpec for this workload.
	PodTemplateSpec() corev1.PodTemplateSpec

	// Select all Workloads of this type with the given selector. May return
	// the current spec if it matches the selector.
	Select(selector map[string]string) ([]TemplateRolloutTarget, error)

	// SelectOwnPods returns a list of pods that are managed by this workload.
	SelectOwnPods(pods []corev1.Pod) ([]corev1.Pod, error)

	// NumReplicas returns the current number of running replicas for this workload.
	NumReplicas() int32

	// PatchNumReplicas modifies the number of replicas for this workload.
	PatchNumReplicas(num int32) error
}

// Deployer is an interface for rollout strategies.
type Deployer interface {
	// Deploy initiates a rollout for a target spec based on the underlying strategy implementation.
	Deploy(cv *cv1.ContainerVersion, version string, spec RolloutTarget) error
}

// SimpleDeployer implements a rollout strategy by patching the target's pod spec with a new version.
type SimpleDeployer struct {
	namespace string

	hp history.Provider

	cs       kubernetes.Interface
	recorder events.Recorder

	opts *conf.Options
}

// NewSimpleDeployer returns a new SimpleDeployer instance, which triggers rollouts
// by patching the target's pod spec with a new version and using the default
// Kubernetes deployment strategy for the workload.
func NewSimpleDeployer(cs kubernetes.Interface, eventRecorder events.Recorder, namespace string,
	options ...func(*conf.Options)) *SimpleDeployer {
	opts := conf.NewOptions()
	for _, opt := range options {
		opt(opts)
	}
	return &SimpleDeployer{
		namespace: namespace,
		cs:        cs,
		recorder:  eventRecorder,
		opts:      opts,
	}
}

// Deploy implements the Deployer interface.
func (sd *SimpleDeployer) Deploy(cv *cv1.ContainerVersion, version string, target RolloutTarget) error {
	log.Printf("Performing simple deployment on %s with version %s", target.Name(), version)

	podSpec := target.PodSpec()

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		for _, c := range podSpec.Containers {
			if c.Name == cv.Spec.Container.Name {
				if updateErr := target.PatchPodSpec(cv, c, version); updateErr != nil {
					log.Printf("Failed to update container version (will retry): version=%v, target=%v, error=%v",
						version, target.Name(), updateErr)

					if updateErr != nil {
						return updateErr
					}
					if sd.opts.UseRollback && target.RollbackAfter() != nil {
						time.Sleep(*target.RollbackAfter())
						if !target.ProgressHealth() {
							currentVersion := strings.SplitAfterN(c.Image, ":", 2)[1]
							return retry.RetryOnConflict(retry.DefaultRetry, func() error {
								if rbErr := target.PatchPodSpec(cv, c, currentVersion); rbErr != nil {
									log.Printf(`Failed to rollback container version (will retry):
										from version=%v, to version=%v, target=%v, error=%v`,
										version, currentVersion, target.Name(), updateErr)
								}
								return nil
							})
						}
					}
					return nil
				}
			}
		}
		return nil
	})
	if retryErr == nil {
		if sd.opts.UseHistory {
			err := sd.hp.Add(sd.namespace, target.Name(), &history.Record{
				Type:    target.Type(),
				Name:    target.Name(),
				Version: version,
				Time:    time.Now(),
			})
			if err != nil {
				sd.opts.Stats.IncCount(fmt.Sprintf("crsyn.%s.history.save.failure", target.Name()))
				sd.recorder.Event(events.Warning, "SaveHistoryFailed", "Failed to record update history")
			}
		}

		log.Printf("Update completed: target=%v", target.Name())
		sd.opts.Stats.IncCount(fmt.Sprintf("crsyn.%s.sync.success", target.Name()))
		sd.recorder.Event(events.Normal, "Success", "Updated completed successfully")
	} else {
		sd.opts.Stats.Event(fmt.Sprintf("%s.sync.failure", target.Name()),
			fmt.Sprintf("Failed to validate image with %s", version), "", "error",
			time.Now().UTC())
		log.Printf("Failed to update container version after maximum retries: version=%v, target=%v, error=%v",
			version, target.Name(), retryErr)
		sd.recorder.Event(events.Warning, "CRSyncFailed", "Failed to deploy the target")
	}
	return nil
}