#!/usr/bin/env bash

# Simple script to generate a timeseries graph from a data file.
# Depencies: gnuplot

# arg 1: path to file with the raw data
# arg 2: path to file to save the graph in (.png will be appended)

gnuplot << EOF
set datafile separator ","

set style data boxes
set boxwidth 0.75
set style fill solid
set grid
unset key
set terminal png size 800,200 transparent truecolor
set output "$2.png"

plot "$1" using 2:xtic(1)
EOF
