package db

import (
	"math"
	"strings"
	"time"

	"github.com/fanaticscripter/EggContractor/coop"
)

type Peeked struct {
	ContractId                 string
	Code                       string
	LastPeekedTime             time.Time
	HasCompleted               bool
	Openings                   int32
	EggsLaid                   float64
	EggsPerHour                float64
	RequiredEggsPerHour        float64
	TimeLeft                   time.Duration
	MaxEarningBonusPercentage  float64
	MeanEarningBonusPercentage float64
}

type Event struct {
	Id            string
	EventType     string
	Multiplier    float64
	Message       string
	FirstSeenTime time.Time
	LastSeenTime  time.Time
	ExpiryTime    time.Time
}

func NewPeeked(c *coop.CoopStatus, peekedAt time.Time) *Peeked {
	var openings int32
	var requiredEggsPerHour float64
	if c.Contract != nil {
		openings = c.Contract.MaxCoopSize - int32(len(c.Members))
		requiredEggsPerHour = c.RequiredEggsPerHour(c.Contract)
	}
	var maxEBP float64
	var sumOoM float64
	var meanEBP float64
	for _, m := range c.Members {
		sumOoM += m.EarningBonusOom
		if m.EarningBonusPercentage() > maxEBP {
			maxEBP = m.EarningBonusPercentage()
		}
	}
	if len(c.Members) > 0 {
		meanEBP = math.Pow(10, sumOoM/float64(len(c.Members))+2)
	}
	return &Peeked{
		ContractId:                 c.ContractId,
		Code:                       c.Code,
		LastPeekedTime:             peekedAt,
		HasCompleted:               c.HasCompleted(),
		Openings:                   openings,
		EggsLaid:                   c.EggsLaid,
		EggsPerHour:                c.EggsPerHour(),
		RequiredEggsPerHour:        requiredEggsPerHour,
		TimeLeft:                   c.DurationUntilProductionDeadline(),
		MaxEarningBonusPercentage:  maxEBP,
		MeanEarningBonusPercentage: meanEBP,
	}
}

func (p *Peeked) HasNoTimeLeft() bool {
	return p.TimeLeft <= 0
}

func (p *Peeked) IsOnTrackToFinish() bool {
	if p.HasCompleted {
		return true
	}
	return p.EggsPerHour >= p.RequiredEggsPerHour
}

func (e *Event) Duration() time.Duration {
	return e.ExpiryTime.Sub(e.FirstSeenTime)
}

func (e *Event) UnhypedMessage() string {
	return strings.ToLower(strings.TrimRight(e.Message, "!"))
}

func (e *Event) HasTimeLeft() bool {
	return e.ExpiryTime.After(time.Now())
}

func (e *Event) TimeLeft() time.Duration {
	return time.Until(e.ExpiryTime)
}
