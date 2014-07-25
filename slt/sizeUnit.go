package slt

func SizeUnit(size int64) (outSize float64, unit string) {
	const tokB = float64(1. / (1024))
	const toMB = float64(1. / (1024*1024))
	const toGB = float64(1. / (1024*1024*1024))
	
	switch {
		case size > (1024*1024*1024):
			return float64(size)*toGB, "GB"
		case size > (1024*1024):
			return float64(size)*toMB, "MB"
		case size > 1024:
			return float64(size)*tokB, "kB"
		default:
			return float64(size), "bytes"
	}
}

