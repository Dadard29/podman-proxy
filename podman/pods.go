package podman

import (
	"fmt"

	"github.com/Dadard29/podman-proxy/models"
	"github.com/containers/libpod/libpod"
)

func (r *PodmanRuntime) GetInfraContainerFromPod(pod *libpod.Pod) (*models.Container, error) {
	var infra *models.Container
	if pod.HasInfraContainer() {
		infraId, err := pod.InfraContainerID()
		if err != nil {
			return nil, err
		}
		infra, err = r.GetContainerFromID(infraId)
		if err != nil {
			return nil, err
		}
	}

	return infra, nil
}

func (r *PodmanRuntime) ListPods() ([]*models.PodmanPod, error) {
	pods, err := r.runtime.GetAllPods()
	if err != nil {
		return nil, err
	}

	out := make([]*models.PodmanPod, 0)
	for _, pod := range pods {

		infra, err := r.GetInfraContainerFromPod(pod)
		if err != nil {
			return nil, err
		}

		podmanPod, err := models.NewPodmanPod(pod, infra)
		if err != nil {
			return nil, err
		}
		out = append(out, podmanPod)
	}

	return out, nil
}

func (r *PodmanRuntime) GetPodFromID(podId string) (*models.PodmanPod, error) {
	pod, err := r.runtime.GetPod(podId)
	if err != nil {
		return nil, err
	}

	infra, err := r.GetInfraContainerFromPod(pod)
	if err != nil {
		return nil, err
	}

	podmanPod, err := models.NewPodmanPod(pod, infra)
	if err != nil {
		return nil, err
	}

	return podmanPod, nil
}

func (r *PodmanRuntime) GetPodFromName(podName string) (*models.PodmanPod, error) {
	pods, err := r.ListPods()
	if err != nil {
		return nil, err
	}

	for _, pod := range pods {
		if pod.Name == podName {
			return pod, nil
		}
	}

	return nil, fmt.Errorf("pod with name %s not found", podName)
}
