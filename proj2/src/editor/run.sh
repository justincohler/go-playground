#=========================================================================
# Sequential
#=========================================================================
cd csv1 
echo "CSV1 (0): " > ../times.txt
(time ../editor csv_file_1.csv) 2>>../times.txt

cd ../csv2
echo "CSV2(0): " >> ../times.txt
(time ../editor csv_file_2.csv) 2>>../times.txt

cd ../csv3
echo "CSV3(0): " >> ../times.txt
(time ../editor csv_file_3.csv) 2>>../times.txt

#=========================================================================
# Parallel
#=========================================================================
cd ../csv1
echo "CSV1 (1): " >> ../times.txt
(time ../editor 1 csv_file_1.csv) 2>>../times.txt
echo "CSV1 (2): " >> ../times.txt
(time ../editor 2 csv_file_1.csv) 2>>../times.txt
echo "CSV1 (4): " >> ../times.txt
(time ../editor 4 csv_file_1.csv) 2>>../times.txt
echo "CSV1 (6): " >> ../times.txt
(time ../editor 6 csv_file_1.csv) 2>>../times.txt
echo "CSV1 (8): " >> ../times.txt
(time ../editor 8 csv_file_1.csv) 2>>../times.txt

cd ../csv2
echo "CSV2(1): " >> ../times.txt
(time ../editor 1 csv_file_2.csv) 2>>../times.txt
echo "CSV2(2): " >> ../times.txt
(time ../editor 2 csv_file_2.csv) 2>>../times.txt
echo "CSV2(4): " >> ../times.txt
(time ../editor 4 csv_file_2.csv) 2>>../times.txt
echo "CSV2(6): " >> ../times.txt
(time ../editor 6 csv_file_2.csv) 2>>../times.txt
echo "CSV2(8): " >> ../times.txt
(time ../editor 8 csv_file_2.csv) 2>>../times.txt

cd ../csv3
echo "CSV3(1): " >> ../times.txt
(time ../editor 1 csv_file_3.csv) 2>>../times.txt
echo "CSV3(2): " >> ../times.txt
(time ../editor 2 csv_file_3.csv) 2>>../times.txt
echo "CSV3(4): " >> ../times.txt
(time ../editor 4 csv_file_3.csv) 2>>../times.txt
echo "CSV3(6): " >> ../times.txt
(time ../editor 6 csv_file_3.csv) 2>>../times.txt
echo "CSV3(8): " >> ../times.txt
(time ../editor 8 csv_file_3.csv) 2>>../times.txt
