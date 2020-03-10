// +build windows

package nsqd

// On Windows, file names cannot contain colons.
//Windows上的名字不能包含冒号
func getBackendName(topicName, channelName string) string {
	// backend names, for uniqueness, automatically include the topic... <topic>;<channel>
	backendName := topicName + ";" + channelName
	return backendName
}
