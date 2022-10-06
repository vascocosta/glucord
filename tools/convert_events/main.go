package main

import (
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

const (
	eventsDiscord = "/home/gluon/var/discord/bots/Hawking/events.csv"
	eventsIRC     = "/home/gluon/var/irc/bots/Schumacher/events.csv"
)

func main() {
	input, err := os.Open(eventsDiscord)
	if err != nil {
		log.Fatalln(err)
	}
	defer input.Close()
	all := dataframe.ReadCSV(input, dataframe.HasHeader(false))
	all.SetNames("Category", "Description", "Session", "Date", "Channel", "Picture", "Mention")
	var channel, mention [100]string
	var rows int
	formula1 := all.Filter(
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[Formula 1]",
		},
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[Formula 2]",
		},
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[Formula 3]",
		},
	)
	rows = formula1.Nrow()
	for i := 0; i != rows; i++ {
		channel[i] = "#formula1"
	}
	formula1 = formula1.Mutate(series.New(channel[:rows], series.String, "Channel"))
	for i := 0; i != rows; i++ {
		mention[i] = "notify"
	}
	formula1 = formula1.Mutate(series.New(mention[:rows], series.String, "Mention"))
	motorsport := all.Filter(
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[Indycar]",
		},
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[IMSA]",
		},
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[NASCAR]",
		},
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[MotoGP]",
		},
	)
	rows = motorsport.Nrow()
	for i := 0; i != rows; i++ {
		channel[i] = "#motorsport"
	}
	motorsport = motorsport.Mutate(series.New(channel[:rows], series.String, "Channel"))
	geeks := all.Filter(
		dataframe.F{
			Colname:    "Category",
			Comparator: series.Eq,
			Comparando: "[Space]",
		},
	)
	rows = geeks.Nrow()
	for i := 0; i != rows; i++ {
		channel[i] = "#geeks"
	}
	geeks = geeks.Mutate(series.New(channel[:rows], series.String, "Channel"))
	combined := formula1.RBind(motorsport).RBind(geeks)
	sorted := combined.Arrange(
		dataframe.Sort("Date"),
	)
	output, err := os.Create(eventsIRC)
	if err != nil {
		log.Fatal(err)
	}
	sorted.WriteCSV(output, dataframe.WriteHeader(false))
}
