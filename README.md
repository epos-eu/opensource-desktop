# EPOS Data Portal Installer

## Introduction

EPOS Data Portal Installer is a user-friendly desktop application designed for the local installation of EPOS environments using Docker or Kubernetes. It simplifies the process by providing a graphical interface for the [EPOS Open Source - Docker installer](https://github.com/epos-eu/opensource-docker) and [EPOS Open Source - Kubernetes installer](https://github.com/epos-eu/opensource-kubernetes) CLI tools.

## Usage

Download the installer or executable for your specific system from the [latest release page](https://github.com/epos-eu/opensource-desktop/releases/latest).

| Operating System | Asset                                      | Description                         |
| ---------------- | ------------------------------------------ | ----------------------------------- |
| Windows          | `epos-data-portal-installer-windows.exe`   | Windows `.exe`                      |
| MacOS            | `epos-data-portal-installer-macos.app.zip` | Application in a `.app` zip         |
| MacOS            | `epos-data-portal-installer-macos.pkg`     | `Pkg` installer for the application |
| Linux            | `epos-data-portal-installer-linux`         | Linux executable                    |

**Note:** Some systems or browsers may flag the download due to a lack of code signing. Future releases aim to address this issue.

**Linux:** At present, there is a challenge with building a universal Linux binary that can run on all distributions. If you encounter issues running the binary, consider building directly from source. See [here](https://github.com/wailsapp/wails/issues/2998) for more information about this particular issue.

## About

The application is built using the [Wails](https://wails.io/) framework, seamlessly combining Go and Vue.js for desktop application development.

- **Frontend:** Vue.js with VueRouter and Vuex for routing and state management.
- **Backend:** Go wraps around Docker and Kubernetes CLI tools, providing a user-friendly interface.

Key Go libraries used:

- [minio/selfudate](https://github.com/minio/selfupdate): Facilitates self-updating functionality.
- [wailsapp/wails](https://github.com/wailsapp/wails): Essential for creating the application and using the Wails Go runtime.
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3): Manages user-created environments.

## Development

Follow the official [Wails guide](https://wails.io/docs/gettingstarted/installation) to set up your development environment.

Once you've set up your environment, verify it by running

```bash
wails doctor
```

This command checks if all necessary dependencies are installed and if the environment is correctly set up.

To run the application in development mode, navigate to the project root directory and execute:

```bash
wails dev
```

This command starts the application in development mode. You can make changes to the code and see the changes reflected in real-time. For more information on Wails application development, refer to the [Wails documentation](https://wails.io/docs/gettingstarted/development).

## Building the Application

To build an executable for your system, run the following command:

```bash
wails build --clean
```

After a successful build, you'll find the system-specific executable in the `build/bin` directory. For more information on building an application with Wails, refer to the [Wails documentation](https://wails.io/docs/gettingstarted/building).

## Contributing

If you want to contribute to a project and make it better, your help is very welcome. Contributing is also a great way to learn more about social coding on Github, new technologies and and their ecosystems and how to make constructive, helpful bug reports, feature requests and the noblest of all contributions: a good, clean pull request.

### How to make a clean pull request

Look for a project's contribution instructions. If there are any, follow them.

- Create a personal fork of the project on Github/GitLab.
- Clone the fork on your local machine. Your remote repo on Github/GitLab is called `origin`.
- Add the original repository as a remote called `upstream`.
- If you created your fork a while ago be sure to pull upstream changes into your local repository.
- Create a new branch to work on! Branch from `develop` if it exists, else from `master` or `main`.
- Implement/fix your feature, comment your code.
- Follow the code style of the project, including indentation.
- If the project has tests run them!
- Write or adapt tests as needed.
- Add or change the documentation as needed.
- Squash your commits into a single commit with git's [interactive rebase](https://help.github.com/articles/interactive-rebase). Create a new branch if necessary.
- Push your branch to your fork on Github/GitLab, the remote `origin`.
- From your fork open a pull request in the correct branch. Target the project's `develop` branch if there is one, else go for `master` or `main`!
- …
- If the maintainer requests further changes just push them to your branch. The PR will be updated automatically.
- Once the pull request is approved and merged you can pull the changes from `upstream` to your local repo and delete your extra branch(es).

And last but not least: Always write your commit messages in the present tense. Your commit message should describe what the commit, when applied, does to the code – not what you did to the code.
