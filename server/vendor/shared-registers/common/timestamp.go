package common

import (
	"shared-registers/common/proto"
)

// FindLargestTimeStamp
// find the largest timeStamp from the inputs, sort on request number and if there is a tie, choose
// the one with largest ClientID
func FindLargestTimeStamp(stamps ...*proto.TimeStamp) *proto.TimeStamp {
	maxReqNum := uint64(0)
	var maxTs *proto.TimeStamp
	// find the largest request number
	for _, s := range stamps {
		if s == nil {
			continue
		}
		if s.GetRequestNumber() > maxReqNum {
			maxReqNum = s.GetRequestNumber()
			maxTs = s
		}
	}
	// choose the one with largest clientID if there is an tie in request number
	for _, s := range stamps {
		if s == nil {
			continue
		}
		if s.GetRequestNumber() == maxReqNum && s.ClientID > maxTs.GetClientID() {
			maxTs = s
		}
	}
	return maxTs
}
