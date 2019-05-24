#=========================================================================
# Sequential
#=========================================================================
cd csv1 
(time ../editor csv_file_1.csv) 2>>times.txt

cd ../csv2
(time ../editor csv_file_2.csv) 2>>times.txt

cd ../csv3
(time ../editor csv_file_3.csv) 2>>times.txt

#=========================================================================
# Parallel
#=========================================================================
cd ../csv1
(time ../editor 1 csv_file_1.csv) 2>>times.txt
(time ../editor 2 csv_file_1.csv) 2>>times.txt
(time ../editor 4 csv_file_1.csv) 2>>times.txt
(time ../editor 6 csv_file_1.csv) 2>>times.txt
(time ../editor 8 csv_file_1.csv) 2>>times.txt

cd ../csv2
(time ../editor 1 csv_file_2.csv) 2>>times.txt
(time ../editor 2 csv_file_2.csv) 2>>times.txt
(time ../editor 4 csv_file_2.csv) 2>>times.txt
(time ../editor 6 csv_file_2.csv) 2>>times.txt
(time ../editor 8 csv_file_2.csv) 2>>times.txt

cd ../csv3
(time ../editor 1 csv_file_3.csv) 2>>times.txt
(time ../editor 2 csv_file_3.csv) 2>>times.txt
(time ../editor 4 csv_file_3.csv) 2>>times.txt
(time ../editor 6 csv_file_3.csv) 2>>times.txt
(time ../editor 8 csv_file_3.csv) 2>>times.txt
