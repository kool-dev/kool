#define ApplicationId "{{B26A0699-CADB-4927-82DB-82842ADB2271}"
#define ApplicationGroup "Firework"
#define ApplicationName "Kool"

#include "environment.iss"

[Setup]
AppId={#ApplicationId}
AppName={#ApplicationName}
AppVersion={#ApplicationVersion}
AppVerName={#ApplicationName} {#ApplicationVersion}
AppPublisher=Firework Web & Mobile LTDA
AppComments=From development to production, our tools bring speed and security to software development teams from different stacks, making their development environments reproducible and easy to set up.
AppCopyright=© 2020 kool.dev - Made by Firework
AppPublisherURL=https://kool.dev
AppSupportURL=https://kool.dev
AppUpdatesURL=https://kool.dev
VersionInfoVersion={#ApplicationVersion}
DefaultDirName={autopf}\{#ApplicationGroup}
DisableWelcomePage=No
DisableDirPage=Yes
DisableProgramGroupPage=Yes
DisableReadyPage=Yes
DefaultGroupName={#ApplicationGroup}
LicenseFile=../LICENSE.md
MinVersion=10.0.10240
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin
SetupIconFile=kool.ico
UninstallDisplayIcon={autopf}\{#ApplicationGroup}\kool.ico
UninstallDisplayName={#ApplicationName}
WizardImageStretch=No
WizardImageFile=kool-setup-small-icon.bmp
WizardSmallImageFile=firework-setup-icon.bmp
ArchitecturesInstallIn64BitMode=x64
ChangesEnvironment=true

[Languages]
Name: en; MessagesFile: "compiler:Default.isl"

[Files]
Source: "..\dist\kool.exe"; DestDir: "{autopf}\{#ApplicationGroup}\bin"; Flags: ignoreversion
Source: "kool.ico"; DestDir: "{autopf}\{#ApplicationGroup}"; Flags: ignoreversion

[Registry]
Root: HKLM; Subkey: "SOFTWARE\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "APPLICATION_NAME"; ValueData: "{autopf}\{#ApplicationGroup}\bin\kool.exe"; Flags: uninsdeletevalue

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
    if CurStep = ssPostInstall
     then EnvAddPath(ExpandConstant('{autopf}') + '\' + ExpandConstant('{#ApplicationGroup}') + '\bin');
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
    if CurUninstallStep = usPostUninstall
    then EnvRemovePath(ExpandConstant('{autopf}') + '\' + ExpandConstant('{#ApplicationGroup}') + '\bin');
end;
