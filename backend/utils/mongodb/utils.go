package mongodb

import (
	"github.com/NeverStopDreamingWang/goi"
)

func GetTime() ISODate {
	return ISODate(goi.GetTime())
}
