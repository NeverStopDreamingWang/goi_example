package mongodb

import (
	"github.com/NeverStopDreamingWang/goi/v3"
)

func GetTime() ISODate {
	return ISODate(goi.GetTime())
}
