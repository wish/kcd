package k8s

import (
	"fmt"
	"strings"

	cv1 "github.com/nearmap/cvmanager/gok8s/apis/custom/v1"
	errs "github.com/nearmap/cvmanager/registry/errs"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
)

const job = "Job"

func (k *K8sProvider) syncJobs(cv *cv1.ContainerVersion, version string, listOpts metav1.ListOptions) error {
	ds, err := k.cs.BatchV1().Jobs(k.namespace).List(listOpts)
	if err != nil {
		k.Recorder.Event(k.Pod, corev1.EventTypeWarning, "CRSyncFailed", "Failed to get dependent jobs")
		return errors.Wrap(err, "failed to read jobs ")
	}
	for _, d := range ds.Items {
		if ci, err := k.checkPodSpec(d.Spec.Template, d.Name, version, cv); err != nil {
			if err == errs.ErrVersionMismatch {
				return k.patchPodSpec(d.Spec.Template, d.Name, ci, cv, func(i int) error {
					_, err := k.cs.BatchV1().Jobs(k.namespace).Patch(d.ObjectMeta.Name, types.StrategicMergePatchType,
						[]byte(fmt.Sprintf(podTemplateSpec, d.Spec.Template.Spec.Containers[i].Name, cv.Spec.ImageRepo, version)))
					return err
				}, func() string { return job })
			} else {
				k.raiseSyncPodErrEvents(err, job, d.Name, cv.Spec.Tag, version)
			}
		}
	}
	return nil
}

func (k *K8sProvider) cvJobs(cv *cv1.ContainerVersion, listOpts metav1.ListOptions) ([]*Resource, error) {
	ds, err := k.cs.BatchV1().Jobs(cv.Namespace).List(listOpts)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch job")
	}
	var cvsList []*Resource
	for _, dd := range ds.Items {
		for _, c := range dd.Spec.Template.Spec.Containers {
			if cv.Spec.Container == c.Name {
				cvsList = append(cvsList, &Resource{
					Namespace: cv.Namespace,
					Name:      dd.Name,
					Type:      job,
					Container: c.Name,
					Version:   strings.SplitAfterN(c.Image, ":", 2)[1],
					CV:        cv.Name,
					Tag:       cv.Spec.Tag,
				})
			}
		}
	}

	return cvsList, nil
}