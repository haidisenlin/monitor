for /f "tokens=5" %%i in ('netstat -aon ^| findstr ":9443"') do (
    if not %%i == 0 (
        taskkill /F /PID  %%i
        goto startProgram
    )
)
:startProgram
start cmd /k D:\\Program\\WSO2\\"API Manager"\\2.5.0\\bin\\wso2server.bat
pause