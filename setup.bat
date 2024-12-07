@echo off
:: ตรวจสอบสิทธิ์ Administrator
>nul 2>&1 (
    net session
) || (
    powershell -Command "Start-Process cmd -ArgumentList '/c %~s0' -Verb runAs"
    exit /b
)

:: ตรวจสอบว่า %~dp0 อยู่ใน PATH แล้วหรือไม่
setlocal enabledelayedexpansion
set "newPath="
set "found=0"
for %%A in ("%PATH:;=" "%") do (
    if /I "%%A" == "%~dp0" (
        set "found=1"
    ) else (
        if defined newPath (
            set "newPath=!newPath!;%%A"
        ) else (
            set "newPath=%%A"
        )
    )
)

:: อัปเดต PATH ถ้ายังไม่มี 
if "%found%" == "0" (
    echo C:\burl not found in PATH. Updating PATH...
    set "newPath=%newPath%;%~dp0"
    setx PATH "%newPath%" /M
    echo PATH updated successfully.
) else (
    echo C:\burl is already in PATH.
)

:: แสดงข้อความว่า PATH ถูกอัปเดต
echo Files installed and PATH checked/updated successfully.
