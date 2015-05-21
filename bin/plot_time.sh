#!/usr/bin/env bash

# Simple script to generate a timeseries graph from a data file.
# Depencies: gnuplot

# arg 1: path to file with the raw data
# arg 2: path to file to save the graph in (.png will be appended)
# arg 3: periodic sampling of data, see http://gnuplot.info/docs_4.2/node121.html

gnuplot << EOF
set datafile separator ","
set xdata time
set timefmt "%s" #time format of input data

set style data lines
set grid
unset key
set terminal png size 800,200 transparent truecolor
set output "$2.png"

plot "$1" every $3 using 1:2 linewidth 2
EOF
