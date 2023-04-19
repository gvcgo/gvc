package vtui

import (
	"runtime"

	"github.com/moqsien/gvc/pkgs/utils"
)

var ConversationIcon string = func() string {
	if runtime.GOOS == utils.Windows {
		return "# "
	}
	return "\uEAC7"
}()

var TokenIcon string = func() string {
	if runtime.GOOS == utils.Windows {
		return "T "
	}
	return "\U000F0C24"
}()

var HelpIcon = func() string {
	if runtime.GOOS == utils.Windows {
		return "? "
	}
	return "\U000F02D6"
}()

var PromptIcon = func() string {
	if runtime.GOOS == utils.Windows {
		return "> "
	}
	return "\ueb33"
}()
