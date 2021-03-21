package installation

import "os"

// This function checks if directories
// are existing and creates them if not
func FolderAutoCreation() {
	folders := [3]string{"config", "error-logs", "web"}

	for _, el := range folders {
		if _, err := os.Stat("./" + el); os.IsNotExist(err) {
			_ = os.Mkdir("./"+el, 0644)
		}
	}
}
