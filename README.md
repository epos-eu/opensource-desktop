# EPOS Data Portal Installer

<img width="1594" alt="Screenshot 2024-04-16 alle 11 54 54" src="https://github.com/epos-eu/opensource-desktop/assets/159156158/87627eda-da2d-4772-8070-190b958a6e6a">

## Introduction

EPOS Data Portal Installer is a user-friendly desktop application designed for the local installation of EPOS
environments using Docker or Kubernetes.
It simplifies the process by providing a graphical interface for
the [EPOS Open Source - Docker installer](https://github.com/epos-eu/opensource-docker)
and [EPOS Open Source - Kubernetes installer](https://github.com/epos-eu/opensource-kubernetes) CLI tools.

## Usage

1. **Download the Installer:**
    - Visit the [latest release page](https://github.com/epos-eu/opensource-desktop/releases/latest) to download the
      appropriate installer for your operating system.

   | Operating System | Asset                                      | Description                         |
                                                                                             |------------------|--------------------------------------------|-------------------------------------|
   | Windows          | `epos-data-portal-installer-windows.exe`   | Windows `.exe`                      |
   | MacOS            | `epos-data-portal-installer-macos.app.zip` | Application in a `.app` zip         |
   | MacOS            | `epos-data-portal-installer-macos.pkg`     | `Pkg` installer for the application |
   | Linux            | `epos-data-portal-installer-linux`         | Linux executable                    |

   **Note:** Some systems or browsers may flag the download due to a lack of code signing.
   Future releases aim to address this issue.

   **Linux Users:** Currently, there are challenges with creating a universal Linux binary compatible with all
   distributions.
   If you encounter issues running the binary, consider building directly from the source.
   For more details, see [this issue](https://github.com/wailsapp/wails/issues/2998).

2. **Run the Installer:**
    - Execute the downloaded installer to launch the application.

3. **Install an Environment:**
    - Click the `Install` button in the Home screen.

4. **Select Environment Type:**
    - Choose between Docker or Kubernetes as the platform to install the environment with and click `Next`.

5. **Configure Environment:**
    - Name your environment, give it a version, and click `Next`.
      (For Kubernetes, you’ll also need to select a context.)

6. **Set Environment Variables:**
    - Edit the environment variables as needed, or leave them at their default values if unsure, and click `Next`.

7. **Start Installation:**
    - Click the `Install` button to initiate the installation process.

8. **View Installed Environment:**
    - After the installation completes, click `See in Installed Environments` to view details about the newly created
      environment.

9. **Open EPOS Data Portal:**
    - In the environment details, click `Open in the browser` on the right of the `Data Portal` row to launch the EPOS
      Data Portal in your default browser.
    - **Note: The Data Portal will work, but it will not show any data, this is because this is just a local
      installation of the system, it has not been populated with data yet.**

10. **Populate the Environment with Data:**
    - To add data, return to the installer and click `Populate Environment`.
      This will open a file picker where you can select a directory containing the `.ttl` files to populate the catalog.
      For more information see [EPOS-DCAT-AP](https://github.com/epos-eu/EPOS-DCAT-AP)

11. **Refresh the Portal:**
    - Once the data population is complete, refresh your browser to view the data in the EPOS Data Portal.

12. **Enjoy!**

## About

The application is built using the [Wails](https://wails.io/) framework, seamlessly combining Go and Vue.js for desktop
application development.

- **Frontend:** Vue.js with VueRouter and Vuex for routing and state management.
- **Backend:** Go wraps around Docker and Kubernetes CLI tools, providing a user-friendly interface.

Key Go libraries used:

- [minio/selfupdate](https://github.com/minio/selfupdate): Facilitates self-updating functionality.
- [wailsapp/wails](https://github.com/wailsapp/wails): Essential for creating the application and using the Wails Go
  runtime.
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3): Manages user-created environments.

## Development

Follow the official [Wails guide](https://wails.io/docs/gettingstarted/installation) to set up your development
environment.

Once you've set up your environment, verify it by running

```bash
wails doctor
```

This command checks if all necessary dependencies are installed and if the environment is correctly set up.

To run the application in development mode, navigate to the project root directory and execute:

```bash
wails dev
```

This command starts the application in development mode.
You can make changes to the code and see the changes reflected in real-time.
For more information on Wails application development, refer to
the [Wails documentation](https://wails.io/docs/gettingstarted/development).

## Building the Application

To build an executable for your system, run the following command:

```bash
wails build --clean
```

After a successful build, you'll find the system-specific executable in the `build/bin` directory.
For more information on building an application with Wails, refer to
the [Wails documentation](https://wails.io/docs/gettingstarted/building).

## Contributing

If you want to contribute to a project and make it better, your help is very welcome.
Contributing is also a great way to learn more about social coding on GitHub, new technologies and their ecosystems and
how to make constructive, helpful bug reports, feature requests and the noblest of all contributions: a good, clean pull
request.

### How to make a clean pull request

Look for a project's contribution instructions. If there are any, follow them.

- Create a personal fork of the project on GitHub/GitLab.
- Clone the fork on your local machine. Your remote repo on GitHub/GitLab is called `origin`.
- Add the original repository as a remote called `upstream`.
- If you created your fork a while ago, be sure to pull upstream changes into your local repository.
- Create a new branch to work on! Branch from `develop` if it exists, else from `master` or `main`.
- Implement/fix your feature, comment your code.
- Follow the code style of the project, including indentation.
- If the project has tests run them!
- Write or adapt tests as needed.
- Add or change the documentation as needed.
- Squash your commits into a single commit with
  git's [interactive rebase](https://help.github.com/articles/interactive-rebase).
  Create a new branch if necessary.
- Push your branch to your fork on GitHub/GitLab, the remote `origin`.
- From your fork, open a pull request in the correct branch.
  Target the project's `develop` branch if there is one, else go for `master` or `main`!
- …
- If the maintainer requests further changes, just push them to your branch. The PR will be updated automatically.
- Once the pull request is approved and merged you can pull the changes from `upstream` to your local repo and delete
  your extra branch(es).

And last but not least: Always write your commit messages in the present tense.
Your commit message should describe what the commit does to the code – not what you did to the code.
