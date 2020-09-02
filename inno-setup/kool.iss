#define ApplicationId "{{B26A0699-CADB-4927-82DB-82842ADB2271}"
#define ApplicationDist "../dist"
#define ApplicationGroup "Firework"
#define ApplicationName "Kool"

#include "environment.iss"

[Setup]
AppId={#ApplicationId}
AppName={#ApplicationName}
AppVerName={#ApplicationName}
AppPublisher=Firework Web & Mobile LTDA
AppComments=From development to production, our tools bring speed and security to software development teams from different stacks, making their development environments reproducible and easy to set up.
AppCopyright=© 2020 kool.dev - Made by Firework
AppPublisherURL=https://kool.dev
AppSupportURL=https://kool.dev
AppUpdatesURL=https://kool.dev
DefaultDirName={autopf}\{#ApplicationGroup}
DisableDirPage=Yes
DisableProgramGroupPage=Yes
DisableReadyPage=Yes
DisableFinishedPage=Yes
DisableWelcomePage=Yes
DefaultGroupName={#ApplicationGroup}
LicenseFile=../LICENSE.md
MinVersion=10.0.10240
Compression=lzma
SolidCompression=yes
PrivilegesRequired=admin
SetupIconFile={#ApplicationName}.ico
UninstallDisplayIcon={autopf}\{#ApplicationGroup}\{#ApplicationName}.ico
UninstallDisplayName={#ApplicationName}
WizardImageStretch=No
WizardSmallImageFile={#ApplicationName}-setup-icon.bmp
ArchitecturesInstallIn64BitMode=x64
ChangesEnvironment=true

[Languages]
Name: pt_BR; MessagesFile: "compiler:Languages\BrazilianPortuguese.isl"
Name: en; MessagesFile: "compiler:Default.isl"

[Files]
Source: "{#ApplicationDist}\{#ApplicationName}.exe"; DestDir: "{autopf}\{#ApplicationGroup}\bin"; Flags: ignoreversion
Source: "{#ApplicationName}.ico"; DestDir: "{autopf}\{#ApplicationGroup}"; Flags: ignoreversion

[Registry]
Root: HKLM; Subkey: "SOFTWARE\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "APPLICATION_NAME"; ValueData: "{autopf}\{#ApplicationGroup}\bin\{#ApplicationName}.exe"; Flags: uninsdeletevalue

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
