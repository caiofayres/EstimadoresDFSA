package Eomlee

import (
	Estimator "SisCom/estimator"
	Utils "SisCom/utils"
	"math"
	"time"
)

type Eomlee struct {
}

func (iv *Eomlee) Estimate(metrics Estimator.EstimateMetrics) Estimator.EstimatorRes {
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

		nextFrameSize, estimatedTags := calculateNextParameters(colisions, tagsReaded, frameSize)
		if metrics.LimitFrameSize {
			frameSize = Estimator.GetLimitedFrameSize(estimatedTags)
		} else {
			frameSize = nextFrameSize
		}
	}
	res.RunTime = time.Since(start)
	res.RunIterations = estimatorLoop
	res.EstimatorName = "Eom Lee"
	//fmt.Println(res)
	return res
}

func calculateNextParameters(colisions, tagsReaded, frameSize int) (int, int) {
	yLast := 2.0
	var yActual float64
	//bLast := math.Inf()
	var bActual float64
	threshold := 0.001
	aux := 0.0
	for (yLast - aux) >= threshold {
		bActual = float64(frameSize) / ((yLast * float64(colisions)) + float64(tagsReaded))
		yActual = (1 - math.Exp(-1.0/bActual)) / (bActual * (1 - (1+(1/bActual))*math.Exp(-1.0/bActual)))
		aux = yLast
		yLast = yActual
	}
	nextFrameSize := yActual * float64(colisions)
	estimatedTags := nextFrameSize / bActual
	return int(math.Ceil(nextFrameSize)), int(math.Ceil(estimatedTags))

}
