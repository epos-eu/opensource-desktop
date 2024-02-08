package installation

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// Returns a new platform installation widget. The homeStackedWidget is used to make the back button go back to the home widget
func NewPlatformGui(homeStackedWidget *widgets.QStackedWidget) *widgets.QStackedWidget {
	// Create a new StackWidget for the platform installation steps
	installationWidget := widgets.NewQStackedWidget(nil)

	platformGui := newInstallationGui(installationWidget)
	mainContentWidget := newMainContentWidget(platformGui)

	installationWidget.AddWidget(platformGui.container)

	// Enable the platform steps button
	platformGui.steps.platform.SetEnabled(true)

	// Add the main content widget to the main content layout
	platformGui.mainContent.mainContent.AddWidget(mainContentWidget, 0, 0)

	// Change the behavior of the back button to go back to the home widget
	platformGui.navigation.back.DisconnectClicked()
	platformGui.navigation.back.ConnectClicked(func(checked bool) {
		// Go back to the home widget and remove this widget from the stacked widget (when going back to the home, the installation is considered cancelled, all the steps are lost)
		homeStackedWidget.RemoveWidget(installationWidget)
		homeStackedWidget.SetCurrentIndex(0)
	})

	return installationWidget
}

// Create a new widget for the main content
func newMainContentWidget(instGui *installationGui) *widgets.QWidget {
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

	// Create a variable to store the selected platform
	var selectedPlatform platform

	// Connect the Toggled signal of each button to a slot that unchecks the other button
	dockerButton.ConnectToggled(func(checked bool) {
		if checked {
			kubernetesButton.SetChecked(false)

			selectedPlatform = docker

			// Enable the next button
			instGui.navigation.next.SetEnabled(true)
		} else {
			instGui.navigation.next.SetEnabled(false)
		}
	})
	kubernetesButton.ConnectToggled(func(checked bool) {
		if checked {
			dockerButton.SetChecked(false)

			selectedPlatform = kubernetes

			// Enable the next button
			instGui.navigation.next.SetEnabled(true)
		} else {
			instGui.navigation.next.SetEnabled(false)
		}
	})

	// Initialize the InstallationState
	installationState := installationState{
		stackedWidget: instGui.rootStackedWidget,
		activeStep:    platformStep,
		platform:      selectedPlatform,
	}

	// Set the behavior of the next button
	instGui.navigation.next.DisconnectClicked()
	instGui.navigation.next.ConnectClicked(func(checked bool) {
		// Set the selected platform in the installation state
		installationState.setPlatform(selectedPlatform)
		// Set the next step in the installation state
		installationState.nextStep(platformStep)
		// Create a new environment step widget
		environmentStepWidget := newEnvironmentSetupGui(&installationState)
		// Add the environment step widget to the stacked widget
		installationState.stackedWidget.AddWidget(environmentStepWidget)
		// Go to the next widget in the stacked widget
		installationState.stackedWidget.SetCurrentIndex(installationState.stackedWidget.CurrentIndex() + 1)
	})

	// Add the buttons to the layout
	layout.AddWidget(dockerButton, 0, 0)
	layout.AddWidget(kubernetesButton, 0, 0)

	return widget
}
