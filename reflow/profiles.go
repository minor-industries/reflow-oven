package reflow

import "time"

var Profile1 = NewScheduleRelativeDurations([]Point{
	{Duration(0 * time.Second), 60},
	{Duration(40 * time.Second), 60}, // preheat the element

	{Duration(25 * time.Second), 100},
	{Duration(90 * time.Second), 150},
	{Duration(30 * time.Second), 183},
	{Duration(60 * time.Second), 235},
	{Duration(5 * time.Second), 235},
	{Duration(25 * time.Second), 183},
	{Duration(30 * time.Second), 25},
	{Duration(30 * time.Second), 25},
})

var Profile2 = NewScheduleRelativeDurations([]Point{
	{Duration(0 * time.Second), 45},
	{Duration(40 * time.Second), 45}, // preheat the element

	{Duration(0 * time.Second), 25},

	{Duration(30 * time.Second), 100},
	{Duration(90 * time.Second), 150},
	{Duration(30 * time.Second), 183},
	{Duration(60 * time.Second), 235},
	{Duration(5 * time.Second), 235},
	{Duration(25 * time.Second), 183},
	{Duration(30 * time.Second), 25},
	{Duration(30 * time.Second), 25},
})
