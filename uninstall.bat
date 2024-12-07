@echo off
:: ตรวจสอบสิทธิ์ Administrator
>nul 2>&1 (
    net session
) || (
    powershell -Command "Start-Process cmd -ArgumentList '/c %~s0' -Verb runAs"
    exit /b
)

:: ลบ ออกจาก PATH
setlocal enabledelayedexpansion
set "newPath="
for %%A in ("%PATH:;=" "%") do (
    if /I "%%A" neq "%~dp0" (
        if defined newPath (
            set "newPath=!newPath!;%%A"
        ) else (
            set "newPath=%%A"
        )
    )
)

:: แสดงข้อความว่า Uninstall เสร็จสมบูรณ์
echo Uninstallation completed successfully.
