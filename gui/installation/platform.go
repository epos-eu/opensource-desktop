package installation

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func Platform() *widgets.QWidget {
	// Create a new QWidget and QGridLayout for the platform widget
	platformWidget := widgets.NewQWidget(nil, core.Qt__Widget)
	platformLayout := widgets.NewQVBoxLayout()

	// Create a new QWidget for the steps and tips layout and set its maximum width
	stepsAndTipsWidget := widgets.NewQWidget(nil, 0)
	stepsAndTipsWidget.SetMaximumWidth(200)

	// Create the steps and tips layout and set it on the stepsAndTipsWidget
	stepsAndTipsLayout := widgets.NewQVBoxLayout()
	stepsAndTipsWidget.SetLayout(stepsAndTipsLayout)

	// Add the steps layout and the tips layout to the steps and tips layout
	stepsAndTipsLayout.AddWidget(stepsWidget(), 0, core.Qt__AlignLeft|core.Qt__AlignTop)
	stepsAndTipsLayout.AddWidget(tipsWidget(), 0, core.Qt__AlignLeft|core.Qt__AlignTop)

	// Create the main layout and add the stepsAndTipsWidget and the main content layout to it
	mainLayout := widgets.NewQHBoxLayout()
	mainLayout.AddWidget(stepsAndTipsWidget, 0, 0)
	mainLayout.AddWidget(mainContentWidget(), 0, 0)

	// Add the main layout and the navigation layout to the platform layout
	platformLayout.AddLayout(mainLayout, 0)
	platformLayout.AddWidget(navigationWidget(), 0, core.Qt__AlignBottom)

	// Set the platform layout on the platform widget
	platformWidget.SetLayout(platformLayout)

	return platformWidget
}

func mainContentWidget() *widgets.QWidget {
	// Create a widget and a layout for the main content
	mainContentWidget := widgets.NewQWidget(nil, 0)
	mainContentLayout := widgets.NewQVBoxLayout()

	// Create a QLabel for the main content
	mainContentLabel := widgets.NewQLabel2("Main content", nil, 0)
	mainContentLayout.AddWidget(mainContentLabel, 0, core.Qt__AlignCenter)

	// Set the layout on the widget
	mainContentWidget.SetLayout(mainContentLayout)

	return mainContentWidget
}

func tipsWidget() *widgets.QWidget {
	// Create a widget and a layout for the tips
	tipsWidget := widgets.NewQWidget(nil, 0)
	tipsLayout := widgets.NewQVBoxLayout()

	// Create a QLabel for the tips
	tipsLabel := widgets.NewQLabel2("Tips", nil, 0)
	tipsLabel.SetStyleSheet("font-size: 16px;")
	tipsLayout.AddWidget(tipsLabel, 0, 0)

	// Set the layout on the widget
	tipsWidget.SetLayout(tipsLayout)

	return tipsWidget
}

func stepsWidget() *widgets.QWidget {
	// Create a layout for the steps
	stepsWidget := widgets.NewQWidget(nil, 0)
	stepsLayout := widgets.NewQVBoxLayout()

	// Create a button for each step
	button1 := widgets.NewQPushButton2("Platform", nil)
	button2 := widgets.NewQPushButton2("Variables", nil)
	button3 := widgets.NewQPushButton2("Launch", nil)

	// TODO: Set the size of the buttons to be the same

	// Add the buttons to the steps layout
	stepsLayout.AddWidget(button1, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	stepsLayout.AddWidget(button2, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)
	stepsLayout.AddWidget(button3, 0, core.Qt__AlignHCenter|core.Qt__AlignTop)

	// Set the steps layout on the steps widget
	stepsWidget.SetLayout(stepsLayout)

	// Calculate the fixed height based on the height of the buttons and the spacing of the layout
	fixedHeight := button1.SizeHint().Height()*3 + stepsLayout.Spacing()*2
	stepsWidget.SetMaximumHeight(fixedHeight)

	return stepsWidget
}

func navigationWidget() *widgets.QWidget {
	// Create a widget and a layout for the navigation buttons
	navigationWidget := widgets.NewQWidget(nil, 0)
	navigationLayout := widgets.NewQHBoxLayout()

	// Create the back button
	backButton := widgets.NewQPushButton2("Back", nil)
	navigationLayout.AddWidget(backButton, 0, core.Qt__AlignLeft)

	// Create the next button
	nextButton := widgets.NewQPushButton2("Next", nil)
	navigationLayout.AddWidget(nextButton, 0, core.Qt__AlignRight)

	// Set the layout on the widget
	navigationWidget.SetLayout(navigationLayout)

	return navigationWidget
}
