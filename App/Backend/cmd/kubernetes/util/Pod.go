package util

import v1 "k8s.io/api/core/v1"

func IsPodMatchingService(pod v1.Pod, service v1.Service) bool {
	selector := service.Spec.Selector

	if len(selector) == 0 {
		return false
	}

	for key, value := range selector {
		if pod.Labels[key] != value {
			return false
		}
	}

	return true
}
