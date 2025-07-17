@echo off
cd ..

:: Menu
echo Select a function to run:
echo 1) Cleanup
echo 2) Full Test
echo 3) Makerandom
echo 4) Render all
set /p choice=Enter your choice (1/2/3/4):

if "%choice%"=="1" (
    echo Running Cleanup...
    del *.jpg
    del *.svg
    del *.lrlogic
    echo Cleanup complete.
) else if "%choice%"=="2" (
    echo Running Full Test...
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
    echo Full Test complete.
) else if "%choice%"=="3" (
    echo Running Makerandom...

    python randomgen.py

    echo Random file generation complete.
) else if "%choice%"=="4" (
    echo Rendering all in directory
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

    echo All files processed.
) else (
    echo Invalid choice. Exiting...
    exit /b 1
)

pause
