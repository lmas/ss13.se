#!/usr/bin/env bash

# Simple script to generate a timeseries graph from a data file.
# Depencies: gnuplot

# arg 1: path to file with the raw data
# arg 2: path to file to save the graph in (.png will be appended)

gnuplot << EOF
set datafile separator ","
set xdata time
set timefmt "%s" #time format of input data

set style data lines
set grid
unset key
set terminal png size 800,200 transparent truecolor
set output "$2.png"
set xtics rotate
set format x "%H:%M\n%d/%m"

plot "$1" using 1:2 linewidth 2
EOF
