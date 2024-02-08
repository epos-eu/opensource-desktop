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
	container   *widgets.QWidget // Contains the main layout with all the other widgets
	steps       *StepsWidget
	tips        *TipsWidget
	navigation  *NavigationWidget
	mainContent *MainContentLayout
}

// Contains the buttons for the different installation steps
type StepsWidget struct {
	widget    *widgets.QWidget
	platform  *widgets.QPushButton
	variables *widgets.QPushButton
	install   *widgets.QPushButton
}

// Contains the tips for the different installation steps
type TipsWidget struct {
	tips   *widgets.QLabel
	widget *widgets.QWidget
}

// Contains the navigation buttons (back and next)
type NavigationWidget struct {
	widget *widgets.QWidget
	back   *widgets.QPushButton
	next   *widgets.QPushButton
}

// Contains the layout for the main content of the installation steps
type MainContentLayout struct {
	container   *widgets.QVBoxLayout // Just a container for the main content (used to add the main content to the window)
	mainContent *widgets.QVBoxLayout // The actual layout for the main content
}

// Function to create the installation widget
func NewInstallationGui(stackedWidget *widgets.QStackedWidget) *InstallationGui {
	// Create a new QWidget and QVBoxLayout for the installation widget
	root := widgets.NewQWidget(nil, core.Qt__Widget)
	rootLayout := widgets.NewQVBoxLayout()

	// Get the components for the gui
	stepsWidget := NewStepsWidget()
	tipsWidget := NewTipsWidget()
	navigationWidget := NewNavigationWidget(stackedWidget)
	mainContentWidget := NewMainContentLayout()

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

	return &InstallationGui{
		container:   root,
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

	// Disable the buttons by default
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

func NewMainContentLayout() *MainContentLayout {
	// Create a layout for the main content
	mainContentLayout := widgets.NewQVBoxLayout()

	// Create a group box
	groupBox := widgets.NewQGroupBox(nil)
	groupBoxLayout := widgets.NewQVBoxLayout()
	groupBox.SetLayout(groupBoxLayout)

	// Add the group box to the main content layout
	mainContentLayout.AddWidget(groupBox, 0, 0)

	// Store both the main content layout and the group box layout:
	// - The main content layout is used to add the main content to the window
	// - The group box layout needs to be populated with the main content of the actual step
	return &MainContentLayout{
		container:   mainContentLayout,
		mainContent: groupBoxLayout,
	}
}
