package Estimator

import "time"

type EstimatorInterface interface {
	Estimate(metrics EstimateMetrics) EstimatorRes
}

type EstimatorRes struct {
	EstimatorName   string
	TagsNumber      int
	TotalSlots      int
	TotalEmptySlots int
	TotalColisions  int
	RunIterations   int
	RunTime         time.Duration
}

type EstimateMetrics struct {
	TagsNumber     int
	FrameSize      int
	LimitFrameSize bool
}

func (om *EstimatorRes) AddToOverall(cm EstimatorRes) {
	om.EstimatorName = cm.EstimatorName
	om.TagsNumber = cm.TagsNumber
	om.TotalColisions += cm.TotalColisions
	om.TotalSlots += cm.TotalSlots
	om.TotalEmptySlots += cm.TotalEmptySlots
	om.RunIterations += cm.RunIterations
	om.RunTime += cm.RunTime
}

func (om *EstimatorRes) GetAverage(iterations int) {
	om.TotalColisions /= iterations
	om.TotalEmptySlots /= iterations
	om.TotalSlots /= iterations
	om.RunIterations /= iterations
	om.RunTime /= time.Duration(iterations)
}

func GetLimitedFrameSize(tagsNumber int) int {
	if tagsNumber >= 1 && tagsNumber <= 5 {
		return 4
	}
	if tagsNumber >= 6 && tagsNumber <= 11 {
		return 8
	}
	if tagsNumber >= 12 && tagsNumber <= 22 {
		return 16
	}
	if tagsNumber >= 23 && tagsNumber <= 44 {
		return 32
	}
	if tagsNumber >= 45 && tagsNumber <= 89 {
		return 64
	}
	if tagsNumber >= 90 && tagsNumber <= 177 {
		return 128
	}
	if tagsNumber >= 178 && tagsNumber <= 355 {
		return 256
	}
	if tagsNumber >= 356 && tagsNumber <= 710 {
		return 512
	}
	if tagsNumber >= 711 && tagsNumber <= 1420 {
		return 1024
	}
	if tagsNumber >= 1421 && tagsNumber <= 2840 {
		return 2048
	}
	if tagsNumber >= 2841 && tagsNumber <= 5680 {
		return 4096
	}
	if tagsNumber >= 5681 && tagsNumber <= 11360 {
		return 8192
	}
	return 8192
}
