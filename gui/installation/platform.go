package installation

import (
	"github.com/therecipe/qt/widgets"
)

func PlatformGui(stackedWidget *widgets.QStackedWidget) *widgets.QWidget {
	platformGui := NewInstallationGui(stackedWidget)

	// Enable the platform steps button
	platformGui.steps.platform.SetEnabled(true)
	return platformGui.base
}
