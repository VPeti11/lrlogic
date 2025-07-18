@echo off
cd ..

:: Menu
echo Select a function to run:
echo 1) Cleanup
echo 2) Full Test
echo 3) Makerandom
echo 4) Render all
echo 5) Timed render
echo 6) Compile LRLogic
echo 7) Compile SVG2LR - Go
set /p choice=Enter your choice (1/2/3/4/5/6/7):

if "%choice%"=="1" (
    echo Running Cleanup...
    del *.jpg
    del *.svg
    del *.lrlogic
    echo Cleanup complete.
) else if "%choice%"=="2" (
    echo Running Full Test...
    :: Capture start time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set start_time=%%a%%b%%c%%d

    del *.svg
    del *.jpg
    set GO111MODULE=on
    go build main.go
    move /Y main lrlogic
    copy .\Tests\*.lrlogic .\

    set /p keep_svg=Do you want to keep SVG files after conversion? (y/n):
    set keep_svg=%keep_svg:~0,1%
    set keep_svg=%keep_svg: =%

    del *.svg
    del *.jpg

    for %%f in (*.lrlogic) do (
        if exist %%f (
            echo Processing %%f...
            if /I "%keep_svg%"=="n" (
                lrlogic.exe --file %%f --nosvg --verbose
            ) else (
                lrlogic.exe --file %%f --verbose
            )
            echo.
        )
    )
    del *.lrlogic

    :: Capture end time and calculate elapsed time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set end_time=%%a%%b%%c%%d
    set /a elapsed_time=%end_time% - %start_time%
    echo Full Test complete.
    echo Elapsed Time: %elapsed_time% seconds.
) else if "%choice%"=="3" (
    echo Running Makerandom...
    python randomgen.py
    echo Random file generation complete.
) else if "%choice%"=="4" (
    echo Rendering all in directory
    :: Capture start time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set start_time=%%a%%b%%c%%d

    set /p keep_svg=Keep SVG files after conversion? (y/n):
    set keep_svg=%keep_svg:~0,1%
    set keep_svg=%keep_svg: =%

    del *.svg
    del *.jpg

    for %%f in (*.lrlogic) do (
        if exist %%f (
            echo Processing %%f...
            if /I "%keep_svg%"=="n" (
                lrlogic.exe --file %%f --nosvg --verbose
            ) else (
                lrlogic.exe --file %%f --verbose
            )
            echo.
        )
    )

    :: Capture end time and calculate elapsed time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set end_time=%%a%%b%%c%%d
    set /a elapsed_time=%end_time% - %start_time%
    echo All files processed.
    echo Elapsed Time: %elapsed_time% seconds.
) else if "%choice%"=="5" (
    echo Enter the file path to render
    set /p file_path=File Path:

    if not exist "%file_path%" (
        echo File not found. Exiting...
        exit /b 1
    )

    set /p keep_svg=Do you want to keep SVG file after conversion? (y/n):
    set keep_svg=%keep_svg:~0,1%
    set keep_svg=%keep_svg: =%

    :: Capture start time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set start_time=%%a%%b%%c%%d

    if /I "%keep_svg%"=="n" (
        lrlogic.exe --file "%file_path%" --nosvg --verbose
    ) else (
        lrlogic.exe --file "%file_path%" --verbose
    )

    :: Capture end time and calculate elapsed time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set end_time=%%a%%b%%c%%d
    set /a elapsed_time=%end_time% - %start_time%
    echo Rendering complete.
    echo Elapsed Time: %elapsed_time% seconds.
) else if "%choice%"=="6" (
    :: Capture start time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set start_time=%%a%%b%%c%%d
    go build main.go
    move /Y main lrlogic
    :: Capture end time and calculate elapsed time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set end_time=%%a%%b%%c%%d
    set /a elapsed_time=%end_time% - %start_time%
    echo Elapsed Time: %elapsed_time% seconds.
) else if "%choice%"=="7" (
    :: Capture start time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set start_time=%%a%%b%%c%%d
    cd ./svg2lrlogic/Go
    go build main.go
    move /Y main svg2lr
    :: Capture end time and calculate elapsed time
    for /f "tokens=1-4 delims=:. " %%a in ("%time%") do set end_time=%%a%%b%%c%%d
    set /a elapsed_time=%end_time% - %start_time%
    echo Elapsed Time: %elapsed_time% seconds.
) else (
    echo Invalid choice. Exiting...
    exit /b 1
)

pause
