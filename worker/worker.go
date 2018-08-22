package worker

type Worker struct {
	Period   float64
	TaskTime float64
	Delay    float64
	Do       func() error
}
