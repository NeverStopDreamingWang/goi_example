package mongodb

import (
	"github.com/NeverStopDreamingWang/goi/v2"
)

func GetTime() ISODate {
	return ISODate(goi.GetTime())
}
