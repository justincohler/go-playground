cd csv1 
echo "0 Threads"
(time ../editor csv_file_1.csv) 2>../times.txt
echo "1 Threads"
(time ../editor 1 csv_file_1.csv) 2>>../times.txt
echo "2 Threads"
(time ../editor 2 csv_file_1.csv) 2>>../times.txt
echo "4 Threads"
(time ../editor 4 csv_file_1.csv) 2>>../times.txt
echo "6 Threads"
(time ../editor 6 csv_file_1.csv) 2>>../times.txt
echo "8 Threads"
(time ../editor 8 csv_file_1.csv) 2>>../times.txt
