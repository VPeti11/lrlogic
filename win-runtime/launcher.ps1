# launcher.ps1

function Show-Menu {
    Write-Host "`nSelect a file to run:`n"
    $options = @(
        @{ Display = "LRLogic"; Value = "lrlogic.exe" },
        @{ Display = "SVG2LR - Go"; Value = "svg2lrlogic.exe" },
        @{ Display = "Random generator"; Value = "lrrandomgen.py" },
        @{ Display = "SVG2LR - Py"; Value = "transformsvg2lr.py" }
    )

    for ($i = 0; $i -lt $options.Count; $i++) {
        Write-Host "$($i+1). $($options[$i].Display)"
    }

    $selection = Read-Host "`nEnter the number of your choice"

    if ($selection -match '^\d+$' -and [int]$selection -ge 1 -and [int]$selection -le $options.Count) {
        return $options[$selection - 1].Value
    } else {
        return $null
    }
}

function Ask-YesNo($message) {
    $response = Read-Host "$message (y/n)"
    return $response -match "^[Yy]"
}

# === Get script's directory ===
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Definition

# === Ask user for working directory, default to Desktop\lrlogic ===
$defaultWorkingDir = Join-Path ([Environment]::GetFolderPath("Desktop")) "lrlogic"
$inputWorkingDir = Read-Host "Enter working directory (default: $defaultWorkingDir)"

$workingDir = if ([string]::IsNullOrWhiteSpace($inputWorkingDir)) { $defaultWorkingDir } else { $inputWorkingDir }

# === Create working directory if not exists ===
if (-not (Test-Path $workingDir)) {
    New-Item -Path $workingDir -ItemType Directory | Out-Null
    Write-Host "Created working directory: $workingDir"
}

# === Copy files from script dir to working dir and track copied files ===
$copiedFiles = @()

# Copy each file individually so we can track them
Get-ChildItem -Path $scriptDir -File | ForEach-Object {
    $sourceFile = $_.FullName
    $destFile = Join-Path $workingDir $_.Name
    Copy-Item -Path $sourceFile -Destination $destFile -Force
    $copiedFiles += $destFile
}

# === Change to working directory ===
Push-Location $workingDir

# === Menu and execution ===
$choice = Show-Menu

try {
    switch ($choice) {
        "lrlogic.exe" {
            $inputFile = Read-Host "Enter name or full path to .lrlogic file"
            if (-not [System.IO.Path]::IsPathRooted($inputFile)) {
                $filepath = Join-Path $scriptDir $inputFile
            } else {
                $filepath = $inputFile
            }

            $args = "--file `"$filepath`""
            if (Ask-YesNo "Add --nosvg?")   { $args += " --nosvg" }
            if (Ask-YesNo "Add --nojpg?")   { $args += " --nojpg" }
            if (Ask-YesNo "Add --verbose?") { $args += " --verbose" }

            Start-Process -FilePath ".\lrlogic.exe" -ArgumentList $args -NoNewWindow -Wait
        }

        "svg2lrlogic.exe" {
            $inputFile = Read-Host "Enter name or full path to .lrlogic file"
            if (-not [System.IO.Path]::IsPathRooted($inputFile)) {
                $filepath = Join-Path $scriptDir $inputFile
            } else {
                $filepath = $inputFile
            }

            $args = "--file `"$filepath`""
            if (Ask-YesNo "Add --verbose?") { $args += " --verbose" }

            Start-Process -FilePath ".\svg2lrlogic.exe" -ArgumentList $args -NoNewWindow -Wait
        }

        "lrrandomgen.py" {
            Start-Process "python" -ArgumentList "lrrandomgen.py" -NoNewWindow -Wait
        }

        "transformsvg2lr.py" {
            Start-Process "python" -ArgumentList "transformsvg2lr.py" -NoNewWindow -Wait
        }

        default {
            Write-Host "`nInvalid selection."
        }
    }
}
finally {
    # === Return to original directory ===
    Pop-Location

    # === Delete only the copied files ===
    foreach ($file in $copiedFiles) {
        if (Test-Path $file) {
            try {
                Remove-Item -Path $file -Force
                Write-Host "Deleted copied file: $file"
            }
            catch {
                Write-Warning "Could not delete file: $file"
            }
        }
    }
}
# I hate Windows <3