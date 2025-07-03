@echo off
REM Simple deploy script for Windows

git add .
git commit -m "deploy"
git push

echo Deployment commands executed. 