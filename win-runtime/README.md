# LRLogic Windows Runtime

## Overview

The Windows distribution of LRLogic-launcher is delivered as a single self-extracting executable built using IExpress. Well its not the best and most modern approach but it works. And thats what matters. This executable bundles the core LRLogic binary compiled from Go, helper Python scripts, a Go helper binary, and a PowerShell-based launcher.

---

## How It Works

1. **Self-extraction:**
   When the user runs the packaged `.exe`, IExpress extracts all bundled files into the Windows temporary directory (`%TEMP%`).

2. **Launcher Execution:**
   After extraction, a launcher executable (`run.exe`) is automatically run. This launcher starts a PowerShell script (`launcher.ps1`) that coordinates execution.

3. **Work Directory Selection:**
   The PowerShell script prompts the user to specify a working directory (defaults to `Desktop\lrlogic`). This directory is created if it doesn’t exist.

4. **Dependency Copying:**
   All extracted files are copied from the temporary folder to the specified working directory. This ensures a clean, isolated environment for the session.

5. **User Menu:**
   The launcher displays a text-based menu for the user to choose which LRLogic-related tool to run, including:

   * The core LRLogic executable
   * The Go-based SVG2LR helper binary
   * The Python random generator script
   * The Python SVG2LR transformation script

6. **Command Arguments:**
   Depending on the tool selected, the script prompts for input files and optional command-line arguments (for example verbose flags, disabling certain features).

7. **Execution:**
   The chosen tool runs in the working directory with the specified parameters. Python scripts are launched via `python` command

8. **Cleanup:**
   After execution completes (whether successful or not), the launcher deletes all copied files from the working directory to leave no residual clutter.

---

## Important Details

* **Extraction Folder:** Files are extracted to `%TEMP%` by IExpress automatically.
* **Working Directory:** User-defined; defaults to a folder named `lrlogic` on the Desktop.
* **Python Requirement:** The system must have Python installed and available in PATH for the Python helper scripts to run. Assuming the program was installed via the included setup file it should have installed it using the `install_win_dependencies.ps1` script
* **PowerShell Version:** Tested on a crappy Windows XP computer from the 2000s running AtlasOS. It worked
* **User Interaction:** The launcher is interactive, running in a PowerShell console window.
* **No elevated privileges required:** The launcher only writes files within user-owned directories.

---

## How it works

```
1. User runs the packaged lrlogiclauncher.EXE
2. IExpress extracts files into %TEMP%.
3. `run.exe` launches `launcher.ps1`.
4. User is prompted to choose or enter a working directory.
5. Files are copied to working directory.
6. User selects which LRLogic tool to run.
7. User enters input files and options.
8. Selected tool runs until completion.
9. Copied files in the working directory are deleted.
10. Launcher exits.
```

---

## Notes and Caveats

* The launcher relies on a **batch-to-exe converted launcher** (`run.exe`) to start PowerShell. I didnt link the batch script to this repo since it only contains 2 lines those being:

      @echo off
      powershell.exe -ExecutionPolicy Bypass -File .\launcher.ps1

* Python must be installed separately and accessible via the system PATH for Python scripts to work. As mentioned before the installer should have taken care of that
* The launcher is designed to keep the working directory clean by deleting copied files after use, but some files might not delete if in use.
* The launcher runs in a console window; no GUI is provided.
* Tested on Windows 10 as mentioned before. Since the software is designed to run on AMD64 systems 32-bit is not supported. The run.exe binary is also a 64-bit only binary. And for the dependencies you need [Chocolatey](https://chocolatey.org/) and its only available on Windows 10 and above

---

