package installation

import (
	"fmt"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// step represents the different installation steps
type step int

const (
	platformStep step = iota
	environmentStep
	variablesStep
	installStep
)

// platform represents the platform choice
type platform int

const (
	docker     platform = iota // 0
	kubernetes                 // 1
)

// installationState represents the state of the installation with the user's choices
type installationState struct {
	stackedWidget *widgets.QStackedWidget // The stacked widget to navigate between the different installation steps
	activeStep    step                    // The last step the user has arrived at
	platform      platform                // The platform choice
	environment   Environment
	// Add more fields here
}

// This function transitions to the next step in the installation process only if called from the last step in the chain
func (i *installationState) nextStep(currentStep step) step {
	// Only transition to the next step if the current step is the active step
	if currentStep == i.activeStep {
		i.activeStep++
	}
	// Return the new active step (to check if the transition was successful)
	return i.activeStep
}

func (i *installationState) setEnvironment(environment Environment) {
	i.environment = environment
}

func (i *installationState) setPlatform(p platform) {
	i.platform = p
}

func (i *installationState) log() {
	fmt.Printf("Current Installation State:\n")
	fmt.Printf("Platform: %v\n", i.platform)
	fmt.Printf("Environment: %v\n", i.environment)
	fmt.Printf("Active Step: %v\n", i.activeStep)
}

type installationGui struct {
	container         *widgets.QWidget        // Contains the main layout with all the other widgets
	rootStackedWidget *widgets.QStackedWidget // The stacked widget to navigate between the different installation steps
	steps             *stepsGui
	tips              *tipsGui
	navigation        *navigationGui
	mainContent       *mainContentLayout
}

// Contains the buttons for the different installation steps
type stepsGui struct {
	widget      *widgets.QWidget
	platform    *widgets.QPushButton
	environment *widgets.QPushButton
	variables   *widgets.QPushButton
	install     *widgets.QPushButton
}

// Contains the tips for the different installation steps
type tipsGui struct {
	tips   *widgets.QLabel
	widget *widgets.QWidget
}

// Contains the navigation buttons (back and next)
type navigationGui struct {
	widget *widgets.QWidget
	back   *widgets.QPushButton
	next   *widgets.QPushButton
}

// Contains the layout for the main content of the installation steps
type mainContentLayout struct {
	container   *widgets.QVBoxLayout // Just a container for the main content (used to add the main content to the window)
	mainContent *widgets.QVBoxLayout // The actual layout for the main content
}

// Function to create the installation widget
func newInstallationGui(stackedWidget *widgets.QStackedWidget) *installationGui {
	// Create a new QWidget and QVBoxLayout for the installation widget
	root := widgets.NewQWidget(nil, 0)
	rootLayout := widgets.NewQVBoxLayout()

	// Get the components for the gui
	stepsWidget := newStepsWidget()
	tipsWidget := newTipsWidget()
	navigationWidget := newNavigationWidget(stackedWidget)
	mainContentWidget := newMainContentLayout()

	// Build the steps and tips layout
	stepsAndTipsLayout := widgets.NewQVBoxLayout()
	stepsAndTipsLayout.AddWidget(stepsWidget.widget, 0, core.Qt__AlignLeft|core.Qt__AlignTop)
	stepsAndTipsLayout.AddWidget(tipsWidget.widget, 1, core.Qt__AlignLeft|core.Qt__AlignTop)
	stepsWidget.widget.SetMaximumWidth(200)
	tipsWidget.widget.SetMaximumWidth(200)

	// Build the main layout
	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddLayout(stepsAndTipsLayout, 0)
	mainLayout.AddLayout(mainContentWidget.container, 1)

	// Add the main layout and the navigation widget to the root layout
	rootLayout.AddLayout(mainLayout, 1)
	rootLayout.AddWidget(navigationWidget.widget, 0, core.Qt__AlignBottom)

	// Set the layout on the root widget
	root.SetLayout(rootLayout)

	return &installationGui{
		container:         root,
		rootStackedWidget: stackedWidget,
		steps:             stepsWidget,
		tips:              tipsWidget,
		navigation:        navigationWidget,
		mainContent:       mainContentWidget,
	}
}

func newTipsWidget() *tipsGui {
	// Create a widget and a layout for the tips
	tipsWidget := widgets.NewQWidget(nil, 0)
	tipsLayout := widgets.NewQVBoxLayout()
	tipsWidget.SetLayout(tipsLayout)

	// Create a QLabel for the tips
	tipsLabel := widgets.NewQLabel2("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras ut lacus volutpat, finibus leo non, aliquam nulla. Duis vel cursus quam. Sed mollis ullamcorper sem, ut dapibus libero auctor nec.", nil, 0)
	// tipsLabel.SetStyleSheet("font-size: 16px;")
	tipsLabel.SetWordWrap(true)

	// Create a group box
	groupBox := widgets.NewQGroupBox2("Info", nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Add the QLabel to the group box layout
	groupBoxLayout.AddWidget(tipsLabel, 0, 0)

	// Add the group box to the tips layout
	tipsLayout.AddWidget(groupBox, 0, 0)

	return &tipsGui{
		tips:   tipsLabel,
		widget: tipsWidget,
	}
}

func newStepsWidget() *stepsGui {
	// Create a layout for the steps
	stepsWidget := widgets.NewQWidget(nil, 0)
	stepsLayout := widgets.NewQVBoxLayout()

	// Create a group box
	groupBox := widgets.NewQGroupBox2("Setup steps", nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Create a button for each step
	platform := widgets.NewQPushButton2("Platform", nil)
	environment := widgets.NewQPushButton2("Environment", nil)
	variables := widgets.NewQPushButton2("Variables", nil)
	launch := widgets.NewQPushButton2("Launch", nil)

	// Disable the buttons by default
	platform.SetEnabled(false)
	environment.SetEnabled(false)
	variables.SetEnabled(false)
	launch.SetEnabled(false)

	// Add the buttons to the group box layout
	groupBoxLayout.AddWidget(platform, 0, 0)
	groupBoxLayout.AddWidget(environment, 0, 0)
	groupBoxLayout.AddWidget(variables, 0, 0)
	groupBoxLayout.AddWidget(launch, 0, 0)

	// Add the group box to the steps layout
	stepsLayout.AddWidget(groupBox, 0, 0)

	// Set the layout on the widget
	stepsWidget.SetLayout(stepsLayout)

	return &stepsGui{
		widget:      stepsWidget,
		platform:    platform,
		environment: environment,
		variables:   variables,
		install:     launch,
	}
}

func newNavigationWidget(stackedWidget *widgets.QStackedWidget) *navigationGui {
	// Create a widget and a layout for the navigation buttons
	navigationWidget := widgets.NewQWidget(nil, 0)
	navigationLayout := widgets.NewQHBoxLayout()

	// Create the back button
	backButton := widgets.NewQPushButton2("Back", nil)
	backButton.ConnectClicked(func(checked bool) {
		// Go back to the previous widget in the stacked widget
		stackedWidget.SetCurrentIndex(stackedWidget.CurrentIndex() - 1)
	})
	navigationLayout.AddWidget(backButton, 0, core.Qt__AlignLeft)

	// Create the next button
	nextButton := widgets.NewQPushButton2("Next", nil)
	navigationLayout.AddWidget(nextButton, 0, core.Qt__AlignRight)

	// Disable the next button by default
	nextButton.SetEnabled(false)

	// Set the layout on the widget
	navigationWidget.SetLayout(navigationLayout)

	return &navigationGui{
		widget: navigationWidget,
		back:   backButton,
		next:   nextButton,
	}
}

func newMainContentLayout() *mainContentLayout {
	// Create a layout for the main content
	layout := widgets.NewQVBoxLayout()

	// Create a group box
	groupBox := widgets.NewQGroupBox(nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Add the group box to the main content layout
	layout.AddWidget(groupBox, 0, 0)

	// Store both the main content layout and the group box layout:
	// - The main content layout is used to add the main content to the window
	// - The group box layout needs to be populated with the main content of the actual step
	return &mainContentLayout{
		container:   layout,
		mainContent: groupBoxLayout,
	}
}
