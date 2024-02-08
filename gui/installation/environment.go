package installation

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func newEnvironmentSetupGui(instState *installationState) *widgets.QWidget {
	// Create a general widget for the step
	instGui := newInstallationGui(instState.stackedWidget)
	environmentStepMainContent := newEnviromentStepMainContentWidget(instGui, instState)

	// Enable the environment steps button
	instGui.steps.environment.SetEnabled(true)

	// Add the main content widget to the main content layout
	instGui.mainContent.mainContent.AddWidget(environmentStepMainContent, 0, 0)

	return instGui.container
}

// Create a new widget for the main content
func newEnviromentStepMainContentWidget(instGui *installationGui, instState *installationState) *widgets.QWidget {
	// Create a new QWidget and QVBoxLayout for the main content widget
	widget := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout()
	widget.SetLayout(mainLayout)

	// Create a QLabel for the title
	title := widgets.NewQLabel2("Environment setup", nil, 0)
	title.SetStyleSheet("font-size: 20px; font-weight: bold;")
	title.SetAlignment(core.Qt__AlignCenter)

	// Add the title to the main layout
	mainLayout.AddWidget(title, 0, 0)

	// Create a QFormLayout for the form
	formLayout := widgets.NewQFormLayout(nil)

	// Create the fields for the form
	nameInput := widgets.NewQLineEdit(nil)
	versionInput := widgets.NewQLineEdit(nil)
	nameInput.SetMinimumSize2(200, 0)
	versionInput.SetMinimumSize2(200, 0)

	// Set a maximum length for the inputs
	nameInput.SetMaxLength(20)
	versionInput.SetMaxLength(20)

	// Add the fields to the form layout
	formLayout.AddRow3("Name:", nameInput)
	formLayout.AddRow3("Version:", versionInput)

	// Create a QRegExpValidator for the version input
	versionValidator := gui.NewQRegExpValidator2(core.NewQRegExp2("^[0-9]+(\\.[0-9]+)*$", core.Qt__CaseInsensitive, core.QRegExp__RegExp2), nil)
	versionInput.SetValidator(versionValidator)

	// Create a QRegExpValidator for the name input
	nameValidator := gui.NewQRegExpValidator2(core.NewQRegExp2("^[a-zA-Z0-9_]+$", core.Qt__CaseInsensitive, core.QRegExp__RegExp2), nil)
	nameInput.SetValidator(nameValidator)

	// Connect the textChanged signal of each QLineEdit to a slot that checks if both text fields are not empty
	nameInput.ConnectTextChanged(func(text string) {
		instGui.navigation.next.SetEnabled(nameInput.Text() != "" && versionInput.Text() != "")
	})
	versionInput.ConnectTextChanged(func(text string) {
		instGui.navigation.next.SetEnabled(nameInput.Text() != "" && versionInput.Text() != "")
	})

	// Set he behavior of the next button
	instGui.navigation.next.DisconnectClicked()
	instGui.navigation.next.ConnectClicked(func(checked bool) {
		// Add the environment to the installation state
		instState.setEnvironment(Environment{
			name:    nameInput.Text(),
			version: versionInput.Text(),
		})
		// Go to the next step
		instState.nextStep(environmentStep)

		// TODO: Go to the next step
	})

	// Add the form layout to the main layout
	mainLayout.AddLayout(formLayout, 0)

	return widget
}

type Environment struct {
	name    string
	version string
}
