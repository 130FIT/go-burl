[Setup]
AppName=Burl
AppVersion=1.0
DefaultDirName={pf}\burl
DefaultGroupName=burl
OutputBaseFilename=burl.installer
Compression=lzma
SolidCompression=yes
; ข้อมูล Publisher
VersionInfoCompany=130FIT
PrivilegesRequired=admin

[Files]
Source: "burl.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "sample\*"; DestDir: "{app}"; Flags: ignoreversion
Source: "runner.sample.json"; DestDir: "{app}"; Flags: ignoreversion
Source: "setup.bat"; DestDir: "{app}"; Flags: ignoreversion
Source: "uninstall.bat"; DestDir: "{app}"; Flags: ignoreversion

[Run]
Filename: "{app}\setup.bat"; Parameters: ""; WorkingDir: "{app}"; Flags: runhidden

[UninstallRun]
Filename: "{app}\uninstall.bat"; Parameters: ""; WorkingDir: "{app}"; Flags: runhidden
