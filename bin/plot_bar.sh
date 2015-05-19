#!/usr/bin/env bash

# Simple script to generate a timeseries graph from a data file.
# Depencies: gnuplot

# arg 1: path to file with the raw data
# arg 2: path to file to save the graph in (.png will be appended)
# arg 3: optional title to show on the graph

gnuplot << EOF
set datafile separator ","

set style data boxes
set boxwidth 0.75
set style fill solid
set grid
unset key
set title "$3" # If a second arg was supplied we show it as a title too
set terminal png size 800,200 transparent truecolor
set output "$2.png"

plot "$1" using 2:xtic(1)
EOF
