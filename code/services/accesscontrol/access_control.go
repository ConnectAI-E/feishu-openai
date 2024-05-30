package accesscontrol

import (
	"start-feishubot/initialization"
	"start-feishubot/utils"
	"sync"
)

var accessCountMap = sync.Map{}
var currentDateFlag = ""

/*
CheckAllowAccessThenIncrement If user has accessed more than 100 times according to accessCountMap, return false.
Otherwise, return true and increase the access count by 1
*/
func CheckAllowAccessThenIncrement(userId *string) bool {

	// Begin a new day, clear the accessCountMap
	currentDateAsString := utils.GetCurrentDateAsString()
	if currentDateFlag != currentDateAsString {
		accessCountMap = sync.Map{}
		currentDateFlag = currentDateAsString
	}

	if CheckAllowAccess(userId) {
		accessedCount, ok := accessCountMap.Load(*userId)
		if !ok {
			accessCountMap.Store(*userId, 1)
		} else {
			accessCountMap.Store(*userId, accessedCount.(int)+1)
		}
		return true
	} else {
		return false
	}
}

func CheckAllowAccess(userId *string) bool {

	if initialization.GetConfig().AccessControlMaxCountPerUserPerDay <= 0 {
		return true
	}

	accessedCount, ok := accessCountMap.Load(*userId)

	if !ok {
		accessCountMap.Store(*userId, 0)
		return true
	}

	// If the user has accessed more than 100 times, return false
	if accessedCount.(int) >= initialization.GetConfig().AccessControlMaxCountPerUserPerDay {
		return false
	}

	// Otherwise, return true
	return true
}

func GetCurrentDateFlag() string {
	return currentDateFlag
}

func GetAccessCountMap() *sync.Map {
	return &accessCountMap
}
