#!/usr/bin/env python3

import math
import sys


def calculate_fuel(race_duration, extra_laps, best_lap, fuel_lap):
    laps = math.ceil(((race_duration * 60) / best_lap)) + extra_laps
    fuel = math.ceil(laps * fuel_lap)
    return fuel

def show_usage():
    print("Usage: !simfuel <race_duration> <extra_laps> <best_lap> <fuel_lap>\n\n")
    print("race_duration - how long the race lasts in minutes")
    print("extra_lap - number of extra laps (paranoia)")
    print("fuel_lap - fuel consumption in litters")

arguments = sys.argv[2].split()
if len(arguments) != 4:
    show_usage()
else:
    try:
        fuel = calculate_fuel(float(arguments[0]), float(arguments[1]), float(arguments[2]), float(arguments[3]))
        print(f"Total fuel: {fuel}L")
    except ValueError:
        print("The best lap is in seconds and the fuel per lap is in litters.")
    except ZeroDivisionError:
        print("Your fastest lap must be greater than zero seconds.")
    except Exception as e:
        print("Unexpected Error", "Please report the message below to the developer:\n\n" + str(e).capitalize() + ".")
