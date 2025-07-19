Set-Location ..

Write-Host "Select a function to run:"
Write-Host "1) Cleanup"
Write-Host "2) Full Test"
Write-Host "3) Makerandom"
Write-Host "4) Render all"
Write-Host "5) Timed render"
Write-Host "6) Compile LRLogic"
Write-Host "7) Compile SVG2LR - Go"
$choice = Read-Host "Enter your choice (1/2/3/4/5/6/7)"

switch ($choice) {
    "1" {
        Write-Host "Running Cleanup..."
        Remove-Item *.jpg, *.svg, *.lrlogic -ErrorAction SilentlyContinue
        Write-Host "Cleanup complete."
    }

    "2" {
        Write-Host "Running Full Test..."
        $startTime = Get-Date

        Remove-Item *.svg, *.jpg -ErrorAction SilentlyContinue
        $env:GO111MODULE = "on"
        go build main.go
        Move-Item -Force main lrlogic
        Copy-Item .\Tests\*.lrlogic .\

        $keep_svg = Read-Host "Do you want to keep SVG files after conversion? (y/n)"
        $keep_svg = $keep_svg.Substring(0,1).ToLower()

        Remove-Item *.svg, *.jpg -ErrorAction SilentlyContinue

        Get-ChildItem -Filter *.lrlogic | ForEach-Object {
            Write-Host "Processing $($_.Name)..."
            if ($keep_svg -eq "n") {
                .\lrlogic.exe --file $_.Name --nosvg --verbose
            } else {
                .\lrlogic.exe --file $_.Name --verbose
            }
            Write-Host
        }

        Remove-Item *.lrlogic -ErrorAction SilentlyContinue

        $elapsed = (Get-Date) - $startTime
        Write-Host "Full Test complete."
        Write-Host "Elapsed Time: $($elapsed.TotalSeconds) seconds."
    }

    "3" {
        Write-Host "Running Makerandom..."
        python randomgen.py
        Write-Host "Random file generation complete."
    }

    "4" {
        Write-Host "Rendering all in directory"
        $startTime = Get-Date

        $keep_svg = Read-Host "Keep SVG files after conversion? (y/n)"
        $keep_svg = $keep_svg.Substring(0,1).ToLower()

        Remove-Item *.svg, *.jpg -ErrorAction SilentlyContinue

        Get-ChildItem -Filter *.lrlogic | ForEach-Object {
            Write-Host "Processing $($_.Name)..."
            if ($keep_svg -eq "n") {
                .\lrlogic.exe --file $_.Name --nosvg --verbose
            } else {
                .\lrlogic.exe --file $_.Name --verbose
            }
            Write-Host
        }

        $elapsed = (Get-Date) - $startTime
        Write-Host "All files processed."
        Write-Host "Elapsed Time: $($elapsed.TotalSeconds) seconds."
    }

    "5" {
        $filePath = Read-Host "Enter the file path to render"
        if (!(Test-Path $filePath)) {
            Write-Host "File not found. Exiting..."
            exit 1
        }

        $keep_svg = Read-Host "Do you want to keep SVG file after conversion? (y/n)"
        $keep_svg = $keep_svg.Substring(0,1).ToLower()

        $startTime = Get-Date

        if ($keep_svg -eq "n") {
            .\lrlogic.exe --file $filePath --nosvg --verbose
        } else {
            .\lrlogic.exe --file $filePath --verbose
        }

        $elapsed = (Get-Date) - $startTime
        Write-Host "Rendering complete."
        Write-Host "Elapsed Time: $($elapsed.TotalSeconds) seconds."
    }

    "6" {
        $startTime = Get-Date
        go build main.go
        Move-Item -Force main lrlogic
        $elapsed = (Get-Date) - $startTime
        Write-Host "Elapsed Time: $($elapsed.TotalSeconds) seconds."
    }

    "7" {
        $startTime = Get-Date
        Push-Location ./svg2lrlogic/Go
        go build main.go
        Move-Item -Force main svg2lr
        Pop-Location
        $elapsed = (Get-Date) - $startTime
        Write-Host "Elapsed Time: $($elapsed.TotalSeconds) seconds."
    }

    Default {
        Write-Host "Invalid choice. Exiting..."
        exit 1
    }
}

Pause
