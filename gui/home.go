package gui

import (
	"github.com/epos-eu/opensource-desktop/gui/installation"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func Home(stackedWidget *widgets.QStackedWidget) *widgets.QWidget {
	// Create a new QWidget and QVBoxLayout for the home widget
	homeWidget := widgets.NewQWidget(nil, 0)
	homeLayout := widgets.NewQVBoxLayout()
	homeWidget.SetLayout(homeLayout)

	// Create a QLabel for the central top section and add it to the home layout
	label := widgets.NewQLabel2("Epos Open Source", nil, 0)
	label.SetAlignment(core.Qt__AlignCenter)
	label.SetStyleSheet("font-size: 20px; font-weight: bold;")
	homeLayout.AddWidget(label, 0, 0)

	// Create a QWidget and QVBoxLayout for the buttons
	buttonsWidget := widgets.NewQWidget(nil, 0)
	buttonsLayout := widgets.NewQVBoxLayout()
	buttonsWidget.SetLayout(buttonsLayout)

	// Create the documentation button
	documentationButton := documentationButton()
	buttonsLayout.AddWidget(documentationButton, 0, 0)

	// Create the Install button
	installButton := installButton(stackedWidget)
	buttonsLayout.AddWidget(installButton, 0, 0)

	// Create the Installed Environments button
	installedEnvironmentsButton := installedEnvironmentsButton(stackedWidget)
	buttonsLayout.AddWidget(installedEnvironmentsButton, 0, 0)

	// Add the buttons widget to the home layout
	homeLayout.AddWidget(buttonsWidget, 0, 0)

	return homeWidget
}

func documentationButton() *widgets.QPushButton {
	button := widgets.NewQPushButton2("Open Documentation", nil)
	button.ConnectClicked(func(checked bool) {
		// Open the documentation in the default web browser
		OpenURL("https://epos-eu.github.io/epos-open-source/doc.html")
	})
	return button
}

func installButton(stackedWidget *widgets.QStackedWidget) *widgets.QPushButton {
	button := widgets.NewQPushButton2("Install", nil)
	button.ConnectClicked(func(checked bool) {
		// Get the widget for the platform selection
		platformWidget := installation.Platform()
		// Add the platform widget to the stacked widget
		stackedWidget.AddWidget(platformWidget)
		// Set the platform widget as the current widget
		stackedWidget.SetCurrentWidget(platformWidget)
	})
	return button
}

func installedEnvironmentsButton(stackedWidget *widgets.QStackedWidget) *widgets.QPushButton {
	button := widgets.NewQPushButton2("Installed Environments", nil)
	button.ConnectClicked(func(checked bool) {
		// TODO: Add the installed environments logic
	})
	return button
}
