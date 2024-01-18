package utils

func HandleError(err error) {
	if err != nil {
		LogError("❌ ERROR", err.Error())
		panic("")
	}
}
