[Code]
const EnvironmentKey = 'Environment';

procedure EnvAddPath(installPath: string);
var
    Paths: string;
begin
    { Retrieve current path (use empty string if entry not exists) }
    if not RegQueryStringValue(HKEY_CURRENT_USER, EnvironmentKey, 'Path', Paths) then
        Paths := '';

    if Paths = '' then
        Paths := installPath + ';'
    else
    begin
        { Skip if string already found in path }
        if Pos(';' + Uppercase(installPath) + ';',  ';' + Uppercase(Paths) + ';') > 0 then exit;
        if Pos(';' + Uppercase(installPath) + '\;', ';' + Uppercase(Paths) + ';') > 0 then exit;

        { Append App Install Path to the end of the path variable }
        Log(Format('Right(Paths, 1): [%s]', [Paths[length(Paths)]]));
        if Paths[length(Paths)] = ';' then
            Paths := Paths + installPath + ';'  { don't double up ';' in env(PATH) }
        else
            Paths := Paths + ';' + installPath + ';' ;
    end;

    { Overwrite (or create if missing) path environment variable }
    if RegWriteStringValue(HKEY_CURRENT_USER, EnvironmentKey, 'Path', Paths)
    then Log(Format('The [%s] added to PATH: [%s]', [installPath, Paths]))
    else Log(Format('Error while adding the [%s] to PATH: [%s]', [installPath, Paths]));
end;

procedure EnvRemovePath(installPath: string);
var
    Paths: string;
    P, Offset, DelimLen: Integer;
begin
    { Skip if registry entry not exists }
    if not RegQueryStringValue(HKEY_CURRENT_USER, EnvironmentKey, 'Path', Paths) then
        exit;

    { Skip if string not found in path }
    DelimLen := 1;     { Length(';') }
    P := Pos(';' + Uppercase(installPath) + ';', ';' + Uppercase(Paths) + ';');
    if P = 0 then
    begin
        { perhaps installPath lives in Paths, but terminated by '\;' }
        DelimLen := 2; { Length('\;') }
        P := Pos(';' + Uppercase(installPath) + '\;', ';' + Uppercase(Paths) + ';');
        if P = 0 then exit;
    end;

    { Decide where to start string subset in Delete() operation. }
    if P = 1 then
        Offset := 0
    else
        Offset := 1;
    { Update path variable }
    Delete(Paths, P - Offset, Length(installPath) + DelimLen);

    { Overwrite path environment variable }
    if RegWriteStringValue(HKEY_CURRENT_USER, EnvironmentKey, 'Path', Paths)
    then Log(Format('The [%s] removed from PATH: [%s]', [installPath, Paths]))
    else Log(Format('Error while removing the [%s] from PATH: [%s]', [installPath, Paths]));
end;
