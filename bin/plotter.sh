#!/usr/bin/env bash

# Simple script to generate a timeseries graph from a data file.
# Depencies: gnuplot

# arg 1: path to file with the raw data
# arg 2: path to file to save the graph in (.png will be appended)
# arg 3: optional title to show on the graph

gnuplot << EOF
set datafile separator ","
set xdata time
set timefmt "%s" #time format of input data

set terminal png size 800,200 transparent truecolor
set output "$2.png"
set format x "%d/%m" # format of output time
set grid
set title "$3" # If a second arg was supplied we show it as a title too
unset key

plot "$1" every 4 using 1:2 with lines linewidth 2
EOF
