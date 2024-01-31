package main

import (
	"log"
	"strconv"
	"strings"
)

type Data struct {
	Key   string
	Value string
	RUSH  bool
	REC   bool
	PASS  bool
}

type Athlete struct {
	Name      string
	C         int
	ATT       int
	PASS_YDS  float64
	PASS_TD   float64
	INT       float64
	SACKS     float64
	CAR       float64
	RUSH_YDS  float64
	RUSH_TD   float64
	RUSH_LONG float64
	REC       float64
	REC_YDS   float64
	REC_TD    float64
	REC_LONG  float64
	TGTS      float64
	LOST      float64
}

func (a *Athlete) SetData(data Data) {
	if data.Key == "C/ATT" {
		result := strings.Split(data.Value, "/")
		for index, val := range result {
			i, err := strconv.Atoi(val)
			if err != nil {
				log.Println(err)
				return
			}

			if index == 0 {
				a.C = i
			}

			if index == 1 {
				a.ATT = i
			}
		}

		return
	}

	value, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		log.Println(err)
		return
	}

	switch data.Key {
	case "YDS":
		if data.RUSH {
			a.RUSH_YDS = value
		} else if data.PASS {
			a.PASS_YDS = value
		} else if data.REC {
			a.REC_YDS = value
		}
	case "TD":
		if data.RUSH {
			a.RUSH_TD = value
		} else if data.PASS {
			a.PASS_TD = value
		} else if data.REC {
			a.REC_TD = value
		}
	case "SACKS":
		a.SACKS = value
	case "CAR":
		a.CAR = value
	case "REC":
		if data.REC {
			a.REC = value
		}
	case "TGTS":
		a.TGTS = value
	case "LOST":
		a.LOST = value
	case "INT":
		a.INT = value
	}

}
