package installation

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func NewPlatformGui(rootStackedWidget *widgets.QStackedWidget) *widgets.QWidget {
	// Create a new StackWidget for the platform installation steps
	installationStackedWidget := widgets.NewQStackedWidget(nil)

	platformGui := NewInstallationGui(installationStackedWidget)

	// Enable the platform steps button
	platformGui.steps.platform.SetEnabled(true)

	// Add the main content widget to the main content layout
	platformGui.mainContent.mainContent.AddWidget(newMainContentWidget(), 0, 0)

	// Change the behavior of the back button to go back to the home widget
	platformGui.navigation.back.DisconnectClicked()
	platformGui.navigation.back.ConnectClicked(func(checked bool) {
		rootStackedWidget.SetCurrentIndex(0)
	})

	return platformGui.container
}

// Create a new widget for the main content
func newMainContentWidget() *widgets.QWidget {
	// Create a new QWidget and QVBoxLayout for the main content widget
	widget := widgets.NewQWidget(nil, 0)
	layout := widgets.NewQHBoxLayout()
	widget.SetLayout(layout)

	// Create the buttons for the platforms
	dockerButton := widgets.NewQPushButton2("Docker", nil)
	kubernetesButton := widgets.NewQPushButton2("Kubernetes", nil)
	// Set the buttons to be checkable
	dockerButton.SetCheckable(true)
	kubernetesButton.SetCheckable(true)
	// Set the icons for the buttons
	dockerButton.SetIcon(gui.NewQIcon5("gui/icons/docker-logo.png"))
	kubernetesButton.SetIcon(gui.NewQIcon5("gui/icons/kubernetes-logo.png"))

	// Set the icon size
	iconSize := core.NewQSize2(128, 128)
	dockerButton.SetIconSize(iconSize)
	kubernetesButton.SetIconSize(iconSize)

	// TODO: Disable the button if the platform is not available on the system
	// kubernetesButton.SetEnabled(false)
	// dockerButton.SetEnabled(false)

	// Set the policy for the size of the buttons to be fixed
	dockerButton.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)
	kubernetesButton.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	// Set the tooltips for the buttons
	dockerButton.SetToolTip("Install on Docker")
	kubernetesButton.SetToolTip("Install on Kubernetes")

	// Connect the Toggled signal of each button to a slot that unchecks the other button
	dockerButton.ConnectToggled(func(checked bool) {
		if checked {
			kubernetesButton.SetChecked(false)
		}
	})
	kubernetesButton.ConnectToggled(func(checked bool) {
		if checked {
			dockerButton.SetChecked(false)
		}
	})

	// Add the buttons to the layout
	layout.AddWidget(dockerButton, 0, 0)
	layout.AddWidget(kubernetesButton, 0, 0)

	return widget
}
