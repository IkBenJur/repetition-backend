package types

import "time"

type LoadPresciptionType int

const (
	FIXED                  LoadPresciptionType = iota
	PERCENTAGE_ONE_REP_MAX LoadPresciptionType = iota
	RPE                    LoadPresciptionType = iota
)

type FixedLoadPrescription struct {
	Id        *int
	Weight    float64
	CreatedAt time.Time
}

type PercentageOneRepMaxLoadPrescription struct {
	Id         *int
	Percentage float64
	CreatedAt  time.Time
}

type RPELoadPrescription struct {
	Id        *int
	RPE       float32
	CreatedAt time.Time
}
