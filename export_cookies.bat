@echo off
echo Exporting cookies from Firefox...
yt-dlp --cookies-from-browser firefox --cookies cookies.txt
echo Cookies exported to cookies.txt
pause