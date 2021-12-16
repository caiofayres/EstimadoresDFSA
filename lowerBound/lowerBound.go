package LowerBound

import (
	Estimator "SisCom/estimator"
	Utils "SisCom/utils"
	"time"
)

type LowerBound struct {
}

func (lb *LowerBound) Estimate(metrics Estimator.EstimateMetrics) Estimator.EstimatorRes {
	res := Estimator.EstimatorRes{}
	res.TagsNumber = metrics.TagsNumber

	estimatorLoop := 0
	tags := metrics.TagsNumber
	frameSize := metrics.FrameSize

	start := time.Now()
	for tags > 0 {
		estimatorLoop++
		colisions, emptySlots, tagsReaded := Utils.SimulateTransmission(tags, frameSize)

		res.TotalColisions += colisions
		res.TotalSlots += frameSize
		res.TotalEmptySlots += emptySlots

		tags -= tagsReaded

		if metrics.LimitFrameSize {
			frameSize = Estimator.GetLimitedFrameSize(tagsReaded + (colisions * 2))
		} else {
			frameSize = colisions * 2
		}
	}
	res.RunTime = time.Since(start)
	res.RunIterations = estimatorLoop
	res.EstimatorName = "Lower Bound"
	return res
}
