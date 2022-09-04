package sftp

func LogError(obj interface{}) {
	if Config.Logger == nil {
		return
	}
	Config.Logger.Error(obj)
}

func LogDebug(obj interface{}) {
	if Config.Logger == nil {
		return
	}
	Config.Logger.Debug(obj)
}
