package gui

import (
	"github.com/therecipe/qt/widgets"
)

// MainWindow is the main window of the application
func MainWindow() {
	app := widgets.NewQApplication(0, nil)
	window := widgets.NewQMainWindow(nil, 0)

	// Create a QStackedWidget
	stackedWidget := widgets.NewQStackedWidget(nil)

	// Add the home widget to the stacked widget
	homeWidget := Home(stackedWidget)
	stackedWidget.AddWidget(homeWidget)

	// Set the home widget as the current widget
	stackedWidget.SetCurrentWidget(homeWidget)

	// Set the stacked widget as the central widget for the main window
	window.SetCentralWidget(stackedWidget)

	window.Show()
	app.Exec()
}
