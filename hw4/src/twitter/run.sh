#!/bin/bash

# 99% ADD/remove
(time ./twitter -s) < 50000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter -s) < 100000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter -s) < 500000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter -s) < 1000000ops_99addremove.txt > out.txt 2>>times.txt

(time ./twitter 1 50000) < 50000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 1 100000) < 100000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 1 500000) < 500000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 1 1000000) < 1000000ops_99addremove.txt > out.txt 2>>times.txt

(time ./twitter 2 25000) < 50000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 2 50000) < 100000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 2 250000) < 500000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 2 500000) < 1000000ops_99addremove.txt > out.txt 2>>times.txt

(time ./twitter 4 12500) < 50000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 4 25000) < 100000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 4 125000) < 500000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 4 250000) < 1000000ops_99addremove.txt > out.txt 2>>times.txt

(time ./twitter 6 8333) < 50000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 6 16667) < 100000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 6 83333) < 500000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 6 166667) < 1000000ops_99addremove.txt > out.txt 2>>times.txt

(time ./twitter 8 6250) < 50000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 8 1250) < 100000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 8 6250) < 500000ops_99addremove.txt > out.txt 2>>times.txt
(time ./twitter 8 12500) < 1000000ops_99addremove.txt > out.txt 2>>times.txt


# 99% CONTAINS
(time ./twitter -s) < 50000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter -s) < 100000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter -s) < 500000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter -s) < 1000000ops_99contains.txt > out.txt 2>>times.txt

(time ./twitter 1 50000) < 50000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 1 100000) < 100000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 1 500000) < 500000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 1 1000000) < 1000000ops_99contains.txt > out.txt 2>>times.txt

(time ./twitter 2 25000) < 50000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 2 50000) < 100000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 2 250000) < 500000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 2 500000) < 1000000ops_99contains.txt > out.txt 2>>times.txt

(time ./twitter 4 12500) < 50000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 4 25000) < 100000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 4 125000) < 500000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 4 250000) < 1000000ops_99contains.txt > out.txt 2>>times.txt

(time ./twitter 6 8333) < 50000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 6 16667) < 100000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 6 83333) < 500000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 6 166667) < 1000000ops_99contains.txt > out.txt 2>>times.txt

(time ./twitter 8 6250) < 50000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 8 1250) < 100000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 8 6250) < 500000ops_99contains.txt > out.txt 2>>times.txt
(time ./twitter 8 12500) < 1000000ops_99contains.txt > out.txt 2>>times.txt

# EVEN
(time ./twitter -s) < 50000ops_even.txt > out.txt 2>>times.txt
(time ./twitter -s) < 100000ops_even.txt > out.txt 2>>times.txt
(time ./twitter -s) < 500000ops_even.txt > out.txt 2>>times.txt
(time ./twitter -s) < 1000000ops_even.txt > out.txt 2>>times.txt

(time ./twitter 1 50000) < 50000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 1 100000) < 100000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 1 500000) < 500000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 1 1000000) < 1000000ops_even.txt > out.txt 2>>times.txt

(time ./twitter 2 25000) < 50000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 2 50000) < 100000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 2 250000) < 500000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 2 500000) < 1000000ops_even.txt > out.txt 2>>times.txt

(time ./twitter 4 12500) < 50000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 4 25000) < 100000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 4 125000) < 500000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 4 250000) < 1000000ops_even.txt > out.txt 2>>times.txt

(time ./twitter 6 8333) < 50000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 6 16667) < 100000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 6 83333) < 500000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 6 166667) < 1000000ops_even.txt > out.txt 2>>times.txt

(time ./twitter 8 6250) < 50000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 8 1250) < 100000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 8 6250) < 500000ops_even.txt > out.txt 2>>times.txt
(time ./twitter 8 12500) < 1000000ops_even.txt > out.txt 2>>times.txt