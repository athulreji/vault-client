package utils

import "signal_client/utils/logger"

func CheckError(err error) {
	if err != nil {
		logger.LogError(err.Error())
		panic(err)
	}
}
