package Ivii

import (
	Estimator "SisCom/estimator"
	Utils "SisCom/utils"
	"math"
	"time"
)

type Ivii struct {
}

func (iv *Ivii) Estimate(metrics Estimator.EstimateMetrics) Estimator.EstimatorRes {
	res := Estimator.EstimatorRes{}
	res.TagsNumber = metrics.TagsNumber
	res.EstimatorName = "IV II"

	tags := metrics.TagsNumber
	initFrameSize := metrics.FrameSize
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

		estimatedTags := calculateNextParameters(colisions, emptySlots, tagsReaded, frameSize, initFrameSize)
		frameSize = Estimator.GetLimitedFrameSize(estimatedTags)
	}
	res.RunIterations = estimatorLoop
	res.RunTime = time.Since(start)
	return res
}

func calculateNextParameters(colisions, emptySlots, tagsReaded, frameSize, initFrameSize int) int {

	if colisions != frameSize {
		return auxiliarEquation(float64(frameSize), float64(colisions), float64(emptySlots), float64(tagsReaded))
	}
	if initFrameSize <= 64 {
		return int(math.Ceil(12.047047047*(float64(frameSize)-1) + 2))
	}
	if initFrameSize >= 128 {
		return int(math.Ceil(6.851851850*(float64(frameSize)-1) + 2))
	}
	return int(math.Ceil(9.497497500*(float64(frameSize)-1) + 2))

}

func auxiliarEquation(frameSize, colisions, emptySlots, tagsReaded float64) int {
	ay := -1.0
	by := 0.0
	i := tagsReaded + (2.0 * colisions)
	estimated := i
	for ay < by {
		t := 1.0 - (1.0 / frameSize)
		a0 := math.Pow(t, estimated)
		a1 := (estimated * a0) / (frameSize * t)
		a2 := 1.0 - (a1 + a0)
		a0 = (a0 * frameSize) - emptySlots
		a1 = (a1 * frameSize) - tagsReaded
		a2 = (a2 * frameSize) - colisions
		a0 *= a0
		a1 *= a1
		a2 *= a2
		by = ay
		ay = math.Sqrt(a0 + a1 + a2)
		if estimated == i {
			by = ay + 1.0
		}
		estimated++
	}
	return int(math.Ceil((estimated - 1) - tagsReaded))
}
