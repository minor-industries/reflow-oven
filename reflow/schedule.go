package reflow

import "time"

type Duration time.Duration

type Point struct {
	Time Duration // since midnight
	Val  float64
}

func (p *Point) T() time.Duration {
	return time.Duration(p.Time)
}

type Schedule []Point

func NewScheduleRelativeDurations(points []Point) Schedule {
	var result Schedule

	cur := time.Duration(0)
	for _, p := range points {
		cur += p.T()
		result = append(result, Point{Duration(cur), p.Val})
	}

	return result
}

func (sc Schedule) Val(t time.Duration) float64 {
	n := len(sc)
	i := n - 1
	for ; i >= 0; i-- {
		p := sc[i]
		if t > time.Duration(p.Time) {
			break
		}
	}

	if i == n-1 {
		return sc[n-1].Val
	}

	cur := sc[i]
	next := sc[i+1]

	ratioNext := float64(t-cur.T()) / float64(next.T()-cur.T())
	ratioCur := 1 - ratioNext

	val := cur.Val*ratioCur + ratioNext*next.Val

	return val
}

func (sc Schedule) PrePend(p0 Point) Schedule {
	result := Schedule{p0}

	for _, p := range sc {
		t := p.T() + p0.T()
		result = append(result, Point{Duration(t), p.Val})
	}

	return result
}
