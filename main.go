package main

import (
	"flag"
	"fmt"
	"log"

	Eo "SisCom/eomlee"
	Estimator "SisCom/estimator"
	Iv "SisCom/ivii"
	Lb "SisCom/lowerBound"
	Sh "SisCom/shoute"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Simulator struct {
	initialTagsNumber int
	incrementTags     int
	maxTags           int
	interationsByStep int
	startFrameSize    int
	limitFrameSize    bool
	lowerBound        bool
	shoute            bool
	eomlee            bool
	ivii              bool
	estimators        []Estimator.EstimatorInterface
}

type SimulationResults struct {
	EomLeeRes     []Estimator.EstimatorRes
	IviiRes       []Estimator.EstimatorRes
	SchouteRes    []Estimator.EstimatorRes
	LowerBoundRes []Estimator.EstimatorRes
}

type CoordinatesResults struct {
	TotalSlots plotter.XYs
	EmptySlots plotter.XYs
	Colisions  plotter.XYs
	Iterations plotter.XYs
	RunTime    plotter.XYs
}

func main() {

	initialTag := flag.Int("start-tags", 0, "initial number of tags to start simulation")
	incrementTags := flag.Int("inc-tags", 0, "number of tags to increment by each step")
	maxTags := flag.Int("max-tags", 0, "number of maximum tags to simulate")

	iterationsByStep := flag.Int("replay-step", 0, "number of iterations by each step")

	startFrameSize := flag.Int("start-frame", 0, "initial value of frame size")
	limitFrameSize := flag.Bool("frame-two", false, "limit frame size to a power of two number")

	lowerBound := flag.Bool("lower-bound", false, "enable lower bound estimator")
	shoute := flag.Bool("shoute", false, "enable shoute estimator")
	eomlee := flag.Bool("eom-lee", false, "enable Eom-Lee estimator")
	ivii := flag.Bool("iv-ii", false, "enable IV-II estimator")
	//Lower Bound, o Shoute, o Eom-Lee e o IV-II;
	flag.Parse()

	simulator := NewSimulator(*initialTag, *incrementTags, *maxTags, *iterationsByStep, *startFrameSize, *limitFrameSize, *lowerBound, *shoute, *eomlee, *ivii)
	err := simulator.ValidateFlags()
	if err != nil {
		panic(err)
	}
	res := simulator.StartSimulation()
	//fmt.Printf("%#v\n", res)

	lowerCoordinates := GenerateCoordinates(res.LowerBoundRes)
	iviiCoordinates := GenerateCoordinates(res.IviiRes)
	schouteCoordinates := GenerateCoordinates(res.SchouteRes)
	eomleeCoordinates := GenerateCoordinates(res.EomLeeRes)

	GenerateColisionsGraph(lowerCoordinates, iviiCoordinates, schouteCoordinates, eomleeCoordinates)
	GenerateEmptySlotsGraph(lowerCoordinates, iviiCoordinates, schouteCoordinates, eomleeCoordinates)
	GenerateTotalSlotsGraph(lowerCoordinates, iviiCoordinates, schouteCoordinates, eomleeCoordinates)
	GenerateRunTimeGraph(lowerCoordinates, iviiCoordinates, schouteCoordinates, eomleeCoordinates)
	GenerateIterationsGraph(lowerCoordinates, iviiCoordinates, schouteCoordinates, eomleeCoordinates)

}
func NewSimulator(initialTag, incrementTags, maxTags, iterationsByStep, startFrameSize int, limitFrameSize, lowerBoud, shoute, eomlee, ivii bool) Simulator {
	return Simulator{
		initialTagsNumber: initialTag,
		incrementTags:     incrementTags,
		maxTags:           maxTags,
		interationsByStep: iterationsByStep,
		startFrameSize:    startFrameSize,
		limitFrameSize:    limitFrameSize,
		lowerBound:        lowerBoud,
		shoute:            shoute,
		eomlee:            eomlee,
		ivii:              ivii,
	}
}

func (s *Simulator) StartSimulation() SimulationResults {
	res := SimulationResults{}
	for i := s.initialTagsNumber; i <= s.maxTags; i += s.incrementTags {
		for estimatorIndex := 0; estimatorIndex < len(s.estimators); estimatorIndex++ {
			overall := Estimator.EstimatorRes{}
			for k := 0; k < s.interationsByStep; k++ {
				metrics := s.GenerateMetrics(i)
				estimator := s.estimators[estimatorIndex]
				res := estimator.Estimate(metrics)
				overall.AddToOverall(res)
			}
			//fmt.Printf("Overall: %#v\n", overall)
			overall.GetAverage(s.interationsByStep)
			//fmt.Printf("Average: %#v\n", overall)
			//fmt.Println("time: ", overall.RunTime)
			switch overall.EstimatorName {
			case "Schoute":
				res.SchouteRes = append(res.SchouteRes, overall)
			case "Lower Bound":
				res.LowerBoundRes = append(res.LowerBoundRes, overall)
			case "IV II":
				res.IviiRes = append(res.IviiRes, overall)
			case "Eom Lee":
				res.EomLeeRes = append(res.EomLeeRes, overall)
			}
		}
		if i == s.maxTags {
			i = s.maxTags + 1
		}
	}
	return res
}

func (s *Simulator) GenerateMetrics(tags int) Estimator.EstimateMetrics {

	return Estimator.EstimateMetrics{
		TagsNumber:     tags,
		FrameSize:      s.startFrameSize,
		LimitFrameSize: s.limitFrameSize,
	}
}

func (s *Simulator) ValidateFlags() error {
	s.estimators = []Estimator.EstimatorInterface{}
	if s.initialTagsNumber <= 0 {
		return fmt.Errorf("initial tags number must be greater than 0")
	}
	if s.incrementTags < 0 {
		s.incrementTags = 0
		log.Println("increment tags number readed is smaller than 0, using 0 instead")
	}
	if s.interationsByStep <= 0 {
		s.interationsByStep = 1
		log.Println("iterations per step readed is smaller than 0, using 1 instead")
	}
	if s.startFrameSize <= 0 {
		return fmt.Errorf("start frame size must be greater than 0")
	}
	if s.maxTags < s.initialTagsNumber {
		s.maxTags = s.initialTagsNumber
		log.Println("number of max tags is smaller than initial tags, using max tags = initial tags")
	}
	if s.maxTags > s.initialTagsNumber && s.incrementTags == 0 {
		s.maxTags = s.initialTagsNumber
		log.Println("number of max tags is greater than initial tags but no increment set, using max tags = initial tags to avoid infite execution")
	}
	if s.eomlee {
		EomleeEstimator := &Eo.Eomlee{}
		s.estimators = append(s.estimators, EomleeEstimator)
	}
	if s.shoute {
		shouteEstimator := &Sh.Shoute{}
		s.estimators = append(s.estimators, shouteEstimator)
	}
	if s.ivii {
		iviiEstimator := &Iv.Ivii{}
		s.estimators = append(s.estimators, iviiEstimator)
	}
	if s.lowerBound {
		lowerBoundEstimator := &Lb.LowerBound{}
		s.estimators = append(s.estimators, lowerBoundEstimator)
	}
	if len(s.estimators) == 0 {
		return fmt.Errorf("no estimator enabled for the simulation")
	}
	return nil
}

func GenerateCoordinates(res []Estimator.EstimatorRes) CoordinatesResults {
	cr := CoordinatesResults{}
	ptsTotalSlots := make(plotter.XYs, len(res))
	ptsEmptySlots := make(plotter.XYs, len(res))
	ptsColisions := make(plotter.XYs, len(res))
	ptsRunTime := make(plotter.XYs, len(res))
	ptsIterations := make(plotter.XYs, len(res))
	for i := 0; i < len(res); i++ {
		ptsTotalSlots[i].Y = float64(res[i].TotalSlots)
		ptsTotalSlots[i].X = float64(res[i].TagsNumber)
		ptsEmptySlots[i].Y = float64(res[i].TotalEmptySlots)
		ptsEmptySlots[i].X = float64(res[i].TagsNumber)
		ptsColisions[i].Y = float64(res[i].TotalColisions)
		ptsColisions[i].X = float64(res[i].TagsNumber)
		ptsRunTime[i].Y = float64(res[i].RunTime / 1000)
		ptsRunTime[i].X = float64(res[i].TagsNumber)
		ptsIterations[i].Y = float64(res[i].RunIterations)
		ptsIterations[i].X = float64(res[i].TagsNumber)
	}
	cr.TotalSlots = ptsTotalSlots
	cr.EmptySlots = ptsEmptySlots
	cr.Colisions = ptsColisions
	cr.RunTime = ptsRunTime
	cr.Iterations = ptsIterations
	return cr
}

func GenerateTotalSlotsGraph(lower, ivii, schoute, eomlee CoordinatesResults) {
	p := plot.New()

	p.X.Label.Text = "Tags"
	p.Y.Label.Text = "Total Slots"
	// Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks with commas.

	err := plotutil.AddLinePoints(p,
		"Lower Bound", lower.TotalSlots,
		"Schoute", schoute.TotalSlots,
		"Eom Lee", eomlee.TotalSlots,
		"IV II", ivii.TotalSlots)
	if err != nil {
		panic(err)
	}

	p.Add(plotter.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "TotalSlots.png"); err != nil {
		panic(err)
	}

}

func GenerateColisionsGraph(lower, ivii, schoute, eomlee CoordinatesResults) {
	p := plot.New()

	p.X.Label.Text = "Tags"
	p.Y.Label.Text = "Total Colisions"
	// Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks with commas.

	err := plotutil.AddLinePoints(p,
		"Lower Bound", lower.Colisions,
		"Schoute", schoute.Colisions,
		"Eom Lee", eomlee.Colisions,
		"IV II", ivii.Colisions)
	if err != nil {
		panic(err)
	}

	p.Add(plotter.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "Colisions.png"); err != nil {
		panic(err)
	}
}

func GenerateEmptySlotsGraph(lower, ivii, schoute, eomlee CoordinatesResults) {
	p := plot.New()

	p.X.Label.Text = "Tags"
	p.Y.Label.Text = "Empty Slots"
	// Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks with commas.

	err := plotutil.AddLinePoints(p,
		"Lower Bound", lower.EmptySlots,
		"Schoute", schoute.EmptySlots,
		"Eom Lee", eomlee.EmptySlots,
		"IV II", ivii.EmptySlots)
	if err != nil {
		panic(err)
	}

	p.Add(plotter.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "EmptySlots.png"); err != nil {
		panic(err)
	}
}

func GenerateRunTimeGraph(lower, ivii, schoute, eomlee CoordinatesResults) {
	p := plot.New()

	p.X.Label.Text = "Tags"
	p.Y.Label.Text = "RunTime(MicroSeconds)"

	// Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks with commas.

	err := plotutil.AddLinePoints(p,
		"Lower Bound", lower.RunTime,
		"Schoute", schoute.RunTime,
		"Eom Lee", eomlee.RunTime,
		"IV II", ivii.RunTime)
	if err != nil {
		panic(err)
	}

	p.Add(plotter.NewGrid())

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "RunTime.png"); err != nil {
		panic(err)
	}
}

func GenerateIterationsGraph(lower, ivii, schoute, eomlee CoordinatesResults) {
	p := plot.New()

	p.X.Label.Text = "Tags"
	p.Y.Label.Text = "Iterations"
	// Use a custom tick marker interface implementation with the Ticks function,
	// that computes the default tick marks and re-labels the major ticks with commas.

	err := plotutil.AddLinePoints(p,
		"Lower Bound", lower.Iterations,
		"Schoute", schoute.Iterations,
		"Eom Lee", eomlee.Iterations,
		"IV II", ivii.Iterations)
	if err != nil {
		panic(err)
	}

	p.Add(plotter.NewGrid())
	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "iterations.png"); err != nil {
		panic(err)
	}
}
