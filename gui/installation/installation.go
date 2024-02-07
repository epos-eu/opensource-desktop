package installation

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// Platform represents the platform choice
type Platform int

const (
	Docker     Platform = iota // 0
	Kubernetes                 // 1
)

// InstallationState represents the state of the installation with the user's choices
type InstallationState struct {
	StackedWidget *widgets.QStackedWidget // The stacked widget to navigate between the different installation steps
	Platform      Platform                // The platform choice
	// Add more fields here
}

type InstallationGui struct {
	base        *widgets.QWidget // Contains the main layout for the installation widget
	steps       *StepsWidget
	tips        *TipsWidget
	navigation  *NavigationWidget
	mainContent *MainContentWidget
}

type StepsWidget struct {
	// Contains the buttons for the different installation steps
	widget    *widgets.QWidget
	platform  *widgets.QPushButton
	variables *widgets.QPushButton
	install   *widgets.QPushButton
}

type TipsWidget struct {
	// Contains the tips for the different installation steps
	tips   *widgets.QLabel
	widget *widgets.QWidget
}

type NavigationWidget struct {
	// Contains the navigation buttons (back and next)
	widget *widgets.QWidget
	back   *widgets.QPushButton
	next   *widgets.QPushButton
}

type MainContentWidget struct {
	// Contains the main content for the different installation steps
	mainContent *widgets.QWidget
}

// Function to create the installation widget
func NewInstallationGui(stackedWidget *widgets.QStackedWidget) *InstallationGui {
	// Create a new QWidget and QGridLayout for the platform widget
	platformWidget := widgets.NewQWidget(nil, core.Qt__Widget)
	platformLayout := widgets.NewQVBoxLayout()

	// Get the components of the gui
	stepsWidget := NewStepsWidget()
	tipsWidget := NewTipsWidget()
	navigationWidget := NewNavigationWidget(stackedWidget)
	mainContentWidget := NewMainContentWidget()

	// Build the steps and tips layout
	stepsAndTipsWidget := widgets.NewQWidget(nil, 0)
	stepsAndTipsWidget.SetMaximumWidth(200)
	stepsAndTipsLayout := widgets.NewQVBoxLayout()
	stepsAndTipsWidget.SetLayout(stepsAndTipsLayout)
	stepsAndTipsLayout.AddWidget(stepsWidget.widget, 0, core.Qt__AlignLeft|core.Qt__AlignTop)
	stepsAndTipsLayout.AddWidget(tipsWidget.widget, 1, core.Qt__AlignLeft|core.Qt__AlignTop)

	// Build the main layout
	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddWidget(stepsAndTipsWidget, 0, 0)
	mainLayout.AddWidget(mainContentWidget.mainContent, 0, 0)

	// Add the main layout and the navigation layout to the platform layout
	platformLayout.AddLayout(mainLayout, 1)
	platformLayout.AddWidget(navigationWidget.widget, 0, core.Qt__AlignBottom)

	// Set the platform layout on the platform widget
	platformWidget.SetLayout(platformLayout)

	return &InstallationGui{
		base:        platformWidget,
		steps:       stepsWidget,
		tips:        tipsWidget,
		navigation:  navigationWidget,
		mainContent: mainContentWidget,
	}
}

func NewTipsWidget() *TipsWidget {
	// Create a widget and a layout for the tips
	tipsWidget := widgets.NewQWidget(nil, 0)
	tipsLayout := widgets.NewQVBoxLayout()
	tipsWidget.SetLayout(tipsLayout)

	// Create a QLabel for the tips
	tipsLabel := widgets.NewQLabel2("Tips", nil, 0)
	// tipsLabel.SetStyleSheet("font-size: 16px;")

	// Create a group box
	groupBox := widgets.NewQGroupBox2("Info", nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Add the QLabel to the group box layout
	groupBoxLayout.AddWidget(tipsLabel, 0, 0)

	// Add the group box to the tips layout
	tipsLayout.AddWidget(groupBox, 0, 0)

	return &TipsWidget{
		tips:   tipsLabel,
		widget: tipsWidget,
	}
}

func NewStepsWidget() *StepsWidget {
	// Create a layout for the steps
	stepsWidget := widgets.NewQWidget(nil, 0)
	stepsLayout := widgets.NewQVBoxLayout()

	// Create a group box
	groupBox := widgets.NewQGroupBox2("Setup steps", nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Create a button for each step
	button1 := widgets.NewQPushButton2("Platform", nil)
	button2 := widgets.NewQPushButton2("Variables", nil)
	button3 := widgets.NewQPushButton2("Launch", nil)

	// TODO: Set the size of the buttons to be the same

	button1.SetEnabled(false)
	button2.SetEnabled(false)
	button3.SetEnabled(false)

	// Add the buttons to the group box layout
	groupBoxLayout.AddWidget(button1, 0, 0)
	groupBoxLayout.AddWidget(button2, 0, 0)
	groupBoxLayout.AddWidget(button3, 0, 0)

	// Add the group box to the steps layout
	stepsLayout.AddWidget(groupBox, 0, 0)

	// Set the layout on the widget
	stepsWidget.SetLayout(stepsLayout)

	return &StepsWidget{
		widget:    stepsWidget,
		platform:  button1,
		variables: button2,
		install:   button3,
	}
}

func NewNavigationWidget(stackedWidget *widgets.QStackedWidget) *NavigationWidget {
	// Create a widget and a layout for the navigation buttons
	navigationWidget := widgets.NewQWidget(nil, 0)
	navigationLayout := widgets.NewQHBoxLayout()

	// Create the back button
	backButton := widgets.NewQPushButton2("Back", nil)
	backButton.ConnectClicked(func(checked bool) {
		// Remove the current widget from the stacked widget
		stackedWidget.RemoveWidget(stackedWidget.CurrentWidget())
		// Go back to the previous widget in the stacked widget
		stackedWidget.SetCurrentIndex(stackedWidget.CurrentIndex() - 1)
	})
	navigationLayout.AddWidget(backButton, 0, core.Qt__AlignLeft)

	// Create the next button
	nextButton := widgets.NewQPushButton2("Next", nil)
	navigationLayout.AddWidget(nextButton, 0, core.Qt__AlignRight)

	// Set the layout on the widget
	navigationWidget.SetLayout(navigationLayout)

	return &NavigationWidget{
		widget: navigationWidget,
		back:   backButton,
		next:   nextButton,
	}
}

func NewMainContentWidget() *MainContentWidget {
	// Create a widget and a layout for the main content
	mainContentWidget := widgets.NewQWidget(nil, 0)
	mainContentLayout := widgets.NewQVBoxLayout()

	// Create a group box
	groupBox := widgets.NewQGroupBox(nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Create a QLabel for the main content
	mainContentLabel := widgets.NewQLabel2("Main content", nil, 0)
	groupBoxLayout.AddWidget(mainContentLabel, 0, core.Qt__AlignCenter)

	// Add the group box to the main content layout
	mainContentLayout.AddWidget(groupBox, 0, 0)

	// Set the layout on the widget
	mainContentWidget.SetLayout(mainContentLayout)

	return &MainContentWidget{
		mainContent: mainContentWidget,
	}
}
