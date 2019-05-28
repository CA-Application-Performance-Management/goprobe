@echo off

SET PRE_PWD=%CD%
SET CONF_WD=%PRE_PWD%\internal\config
SET LOG_WD=%PRE_PWD%\internal\logger
SET MET_WD=%PRE_PWD%\internal\metric
SET UTL_WD=%PRE_PWD%\internal\utils


FOR %%p IN (%PRE_PWD% %CONF_WD% %LOG_WD% %MET_WD% %UTL_WD%) DO (
	echo %%p
	cd %%p
	go build
	IF %ERRORLEVEL% == 0 (
		echo "Build" %%p "Finished"
	) ELSE (
		echo "Build" %%p "Failed"
		EXIT /B 1
	)
	
	go test
	IF %ERRORLEVEL% == 0 (
		echo "Test" %%p "Finished"
	) ELSE (
		echo "Test" %%p "Failed"
		EXIT /B 1
	)
)

