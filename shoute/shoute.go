package Shoute

import (
	Estimator "SisCom/estimator"
	Utils "SisCom/utils"
	"math"
	"time"
)

type Shoute struct {
}

func (iv *Shoute) Estimate(metrics Estimator.EstimateMetrics) Estimator.EstimatorRes {
	res := Estimator.EstimatorRes{}
	res.TagsNumber = metrics.TagsNumber

	tags := metrics.TagsNumber
	frameSize := metrics.FrameSize
	estimatorLoop := 0
	start := time.Now()
	for tags > 0 {
		estimatorLoop++
		colisions, emptySlots, tagsReaded := Utils.SimulateTransmission(tags, frameSize)

		res.TotalColisions += colisions
		res.TotalSlots += frameSize
		res.TotalEmptySlots += emptySlots

		tags -= tagsReaded

		if metrics.LimitFrameSize {
			frameSize = Estimator.GetLimitedFrameSize(tagsReaded + int(math.Ceil((float64(colisions) * 2.39))))
		} else {
			frameSize = int(math.Ceil((float64(colisions) * 2.39)))
		}

	}

	res.RunTime = time.Since(start)
	res.RunIterations = estimatorLoop
	res.EstimatorName = "Schoute"
	return res
}
