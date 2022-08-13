#!/bin/bash

EVENT='2022 Hungarian GP'
RACE_STRING='1117/hungary'
TITLE='Race Results'

case $2 in
	fp*)
		if [ $2 = "fp1" ]; then
			SESSION_STRING='practice-1'
			TITLE='Free Practice 1 Results'
		elif [ $2 = "fp2" ]; then
			SESSION_STRING='practice-2'
			TITLE='Free Practice 2 Results'
		else
			SESSION_STRING='practice-3'
			TITLE='Free Practice 3 Results'
		fi
		RESULTS=`lynx -dump "https://www.formula1.com/en/results.html/2022/races/$RACE_STRING/$SESSION_STRING.html" | \
                        grep '.\{10,100\}' | \
                        sed -e 's/Red Bull Racing RBPT/0/g' \
                            -e 's/Alpine Renault/0/g' \
                            -e 's/Alfa Romeo Ferrari/0/g' \
                            -e 's/AlphaTauri RBPT/0/g' \
                            -e 's/Haas Ferrari/0/g' \
                            -e 's/McLaren Mercedes/0/g' \
                            -e 's/Aston Martin Aramco Mercedes/0/g' \
                            -e 's/Williams Mercedes/0/g' | \
                        grep -A 20 -i 'pos no' | \
			grep '[1-9]:[1-9]' | \
                        awk '{print "**"$1".** " $5 " " $7}' | \
                        tail -n +1`
		;;
	q*)
		SESSION_STRING='qualifying'
		TITLE='Qualifying Results'
		RESULTS=`lynx -dump "https://www.formula1.com/en/results.html/2022/races/$RACE_STRING/$SESSION_STRING.html" | \
		        grep '.\{10,100\}' | \
		        sed -e 's/Red Bull Racing RBPT/0/g' \
		            -e 's/Alpine Renault/0/g' \
		            -e 's/Alfa Romeo Ferrari/0/g' \
		            -e 's/AlphaTauri RBPT/0/g' \
		            -e 's/Haas Ferrari/0/g' \
		            -e 's/McLaren Mercedes/0/g' \
		            -e 's/Aston Martin Aramco Mercedes/0/g' \
		            -e 's/Williams Mercedes/0/g' | \
		        grep -A 20 -i 'pos no' | \
		        awk '{print "**"$1".** " $5 " " $9}' | \
		        tail -n +2`
		;;
	r*)
		SESSION_STRING='race-result'
		TITLE='Race Results'
		RESULTS=`lynx -dump "https://www.formula1.com/en/results.html/2022/races/$RACE_STRING/$SESSION_STRING.html" | \
			grep '.\{10,100\}' | \
			grep -A 20 -i 'pos no' | \
			awk '{print "**"$1".** " $5 " " $--NF}' | \
			tail -n +2`
		;;
	s*)
		SESSION_STRING='sprint-results'
		TITLE='Sprint Race Results'
                RESULTS=`lynx -dump "https://www.formula1.com/en/results.html/2022/races/$RACE_STRING/$SESSION_STRING.html" | \
			grep '.\{10,100\}' | \
			grep -A 20 -i 'pos no' | \
			awk '{print "**"$1".** " $5 " " $--NF}' | \
			tail -n +2`
		;;
	*)
		echo "Usage: !results <fp1|fp2|fp3|qualifying|race>"
		exit
		;;
esac

echo -e "**$EVENT $TITLE:**\n"

if [ ${#RESULTS} -gt 0 ]
then
	echo "$RESULTS"
else
	echo "No results available yet. Wait a bit more and try again please."
fi
