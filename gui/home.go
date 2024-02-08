package gui

import (
	"github.com/epos-eu/opensource-desktop/gui/installation"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

func NewHomeWidget(stackedWidget *widgets.QStackedWidget) *widgets.QWidget {
	// Create a new QWidget and QVBoxLayout for the home widget
	homeWidget := widgets.NewQWidget(nil, 0)
	homeLayout := widgets.NewQVBoxLayout()
	homeWidget.SetLayout(homeLayout)

	// Create a central top widget and add it to the home layout
	centralTopWidget := newCentralTopWidget()
	homeLayout.AddWidget(centralTopWidget, 1, core.Qt__AlignHCenter)

	// Create a QWidget and QVBoxLayout for the buttons
	buttonsWidget := widgets.NewQWidget(nil, 0)
	buttonsLayout := widgets.NewQVBoxLayout()
	buttonsWidget.SetLayout(buttonsLayout)
	// Set the dimensions of the buttons widget
	buttonsWidget.SetFixedWidth(200)

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
	homeLayout.AddWidget(buttonsWidget, 0, core.Qt__AlignHCenter)

	return homeWidget
}

func newCentralTopWidget() *widgets.QWidget {
	// Create a Widget for the central top section
	centralTopWidget := widgets.NewQWidget(nil, 0)
	centralTopLayout := widgets.NewQVBoxLayout()
	// Create a QLabel for the titleWidget and add it to the home layout
	titleWidget := widgets.NewQLabel2("Epos Open Source", nil, 0)
	titleWidget.SetAlignment(core.Qt__AlignCenter)
	titleWidget.SetStyleSheet("font-size: 20px; font-weight: bold;")
	centralTopLayout.AddWidget(titleWidget, 0, core.Qt__AlignHCenter)

	// Create a QLabel for the text and add it to the home layout
	textWidget := widgets.NewQLabel2("EPOS Data portal is an open source, service-based data integration system. It is based on a microservices architecture. The current web page includes general documentation about EPOS, technical documentation on the overall functioning of the architecture, and source code for each microservice. The package will also provide multi-platform executables to allow for full installation locally or on one or more remote servers using Docker and Kubernetes. The version is released under the GPLv3 license.", nil, 0)
	textWidget.SetAlignment(core.Qt__AlignCenter)
	textWidget.SetWordWrap(true)
	centralTopLayout.AddWidget(textWidget, 1, core.Qt__AlignHCenter|core.Qt__AlignTop)

	centralTopWidget.SetLayout(centralTopLayout)
	return centralTopWidget
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
		platformWidget := installation.PlatformGui(stackedWidget)
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
