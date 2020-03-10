// +build !windows

package nsqd

//linux上的名字可以包含冒号
func getBackendName(topicName, channelName string) string {
	// backend names, for uniqueness, automatically include the topic... <topic>:<channel>
	backendName := topicName + ":" + channelName
	return backendName
}
