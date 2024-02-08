package installation

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func newVariablesGui(instState *installationState) *widgets.QWidget {
	// Create a general widget for the step
	instGui := newInstallationGui(instState.stackedWidget)
	variablesStepMainContent := newVariablesStepMainContentWidget(instGui, instState)

	// Enable the variables steps button
	instGui.steps.variables.SetEnabled(true)

	// Add the main content widget to the main content layout
	instGui.mainContent.mainContent.AddWidget(variablesStepMainContent, 0, 0)

	return instGui.container
}

func newVariablesStepMainContentWidget(instGui *installationGui, instState *installationState) *widgets.QWidget {
	// Create a new QWidget and QVBoxLayout for the main content widget
	widget := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout()
	widget.SetLayout(mainLayout)

	// Create a QLabel for the title
	title := widgets.NewQLabel2("Environment Variables", nil, 0)
	title.SetStyleSheet("font-size: 18px; font-weight: bold;")
	title.SetAlignment(core.Qt__AlignCenter)
	mainLayout.AddWidget(title, 0, 0)

	// Create a QWidget for the forms and a QVBoxLayout for its layout
	formsWidget := widgets.NewQWidget(nil, 0)
	formsLayout := widgets.NewQVBoxLayout()
	formsWidget.SetLayout(formsLayout)

	// TODO: Do this for real (use the env file?)
	form := newFormLayout("General variables", map[string]string{
		"Var1": "Default1",
		"Var2": "Default2",
		"Var3": "Default3",
		"Var4": "Default4",
	})
	formsLayout.AddWidget(form.widget, 0, core.Qt__AlignHCenter)
	form2 := newFormLayout("Container 1", map[string]string{
		"Var1": "Default1",
		"Var2": "Default2",
		"Var3": "Default3",
		"Var4": "Default4",
	})
	formsLayout.AddWidget(form2.widget, 0, core.Qt__AlignHCenter)
	form3 := newFormLayout("Container 2", map[string]string{
		"Var1": "Default1",
		"Var2": "Default2",
		"Var3": "Default3",
		"Var4": "Default4",
	})
	formsLayout.AddWidget(form3.widget, 0, core.Qt__AlignHCenter)

	// Create a QScrollArea and set the forms widget as its widget
	scrollArea := widgets.NewQScrollArea(nil)
	scrollArea.SetWidget(formsWidget)
	scrollArea.SetWidgetResizable(true)

	// Add the scroll area to the main layout
	mainLayout.AddWidget(scrollArea, 0, 0)

	return widget
}

// formWidget is a struct that holds a QWidget and a map of QLineEdits
type formWidget struct {
	widget     *widgets.QWidget
	textFields map[string]*widgets.QLineEdit
}

func newFormLayout(title string, fields map[string]string) formWidget {
	widget := widgets.NewQWidget(nil, 0)
	mainLayout := widgets.NewQVBoxLayout()
	widget.SetLayout(mainLayout)
	widget.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	// Create a QLabel for the title
	titleWidget := widgets.NewQLabel2(title, nil, 0)
	titleWidget.SetStyleSheet("font-size: 18px; font-weight: bold;")

	// Add the title to the main layout
	mainLayout.AddWidget(titleWidget, 0, core.Qt__AlignHCenter)

	// Create a QFormLayout for the form
	formLayout := widgets.NewQFormLayout(nil)

	// Create the fields for the form
	textFields := make(map[string]*widgets.QLineEdit)
	for fieldName, defaultValue := range fields {
		// Create a QLineEdit for the text field
		textField := widgets.NewQLineEdit(nil)
		textField.SetMinimumWidth(250)
		textField.SetText(defaultValue) // Set the default value

		// Add the label and text field to the form layout
		formLayout.AddRow3(fieldName+":", textField)

		// Add the text field to the map
		textFields[fieldName] = textField
	}

	// Add the form layout to the main layout
	mainLayout.AddLayout(formLayout, 0)

	return formWidget{
		widget:     widget,
		textFields: textFields,
	}
}
