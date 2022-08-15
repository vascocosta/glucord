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
    print("best_lap - best lap time in seconds")
    print("fuel_lap - fuel consumption in litters")


if len(sys.argv) != 6:
    show_usage()
    exit()
else:
    try:
        fuel = calculate_fuel(float(sys.argv[2]), float(sys.argv[3]), float(sys.argv[4]), float(sys.argv[5]))
        print(f"Total fuel: {fuel}L")
    except ValueError:
        print("Make sure you provide numbers as arguments instead of other characters.")
        print("The race duration is in minutes, the extra laps is an integer.")
        print("The best lap is in seconds and the fuel per lap is in litters.")
    except ZeroDivisionError:
        print("Your fastest lap must be greater than zero seconds.")
    except Exception as e:
        print("Unexpected error, please report the message below to the developer:\n\n" + str(e).capitalize() + ".")
