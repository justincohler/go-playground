echo "Starting Single Thread..." > times.txt
(time ./nbody.exe -bodies=10 -steps=10000 -daysPerStep=5 -threads=1) 2>>times.txt
(time ./nbody.exe -bodies=10 -steps=100000 -daysPerStep=5 -threads=1) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=10000 -daysPerStep=5 -threads=1) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=100000 -daysPerStep=5 -threads=1) 2>>times.txt

echo "Starting Two Thread..." >> times.txt
(time ./nbody.exe -bodies=10 -steps=10000 -daysPerStep=5 -threads=2) 2>>times.txt
(time ./nbody.exe -bodies=10 -steps=100000 -daysPerStep=5 -threads=2) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=10000 -daysPerStep=5 -threads=2) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=100000 -daysPerStep=5 -threads=2) 2>>times.txt

echo "Starting Two Thread..." >> times.txt
(time ./nbody.exe -bodies=10 -steps=10000 -daysPerStep=5 -threads=4) 2>>times.txt
(time ./nbody.exe -bodies=10 -steps=100000 -daysPerStep=5 -threads=4) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=10000 -daysPerStep=5 -threads=4) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=100000 -daysPerStep=5 -threads=4) 2>>times.txt

echo "Starting Two Thread..." >> times.txt
(time ./nbody.exe -bodies=10 -steps=10000 -daysPerStep=5 -threads=8) 2>>times.txt
(time ./nbody.exe -bodies=10 -steps=100000 -daysPerStep=5 -threads=8) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=10000 -daysPerStep=5 -threads=8) 2>>times.txt
(time ./nbody.exe -bodies=50 -steps=100000 -daysPerStep=5 -threads=8) 2>>times.txt