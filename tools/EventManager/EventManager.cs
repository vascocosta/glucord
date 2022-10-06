using Microsoft.Data.Analysis;

namespace EventManager
{
    public static class Resolver {
        public static IDictionary<string, string[]> categories = new Dictionary<string, string[]>
        {
            {"[Formula1]", new string[]
            {
                "665554362570899476",
                "https://logodownload.org/wp-content/uploads/2016/11/formula-1-logo-1-1.png",
                "<@&1005570005682901133>"
            }},
            {"[Formula2]", new string[]
            {
                "665554362570899476",
                "https://upload.wikimedia.org/wikipedia/en/thumb/1/1f/Formula_2_logo.svg/1920px-Formula_2_logo.svg.png",
                ""
            }},
            {"[Formula3]", new string[]
            {
                "665554362570899476",
                "https://upload.wikimedia.org/wikipedia/commons/5/5b/FIA_F3_Championship_logo.png",
                ""
            }},
            {"[IndyCar]", new string[]
            {
                "665554362570899476",
                "https://upload.wikimedia.org/wikipedia/en/thumb/b/bb/INDYCAR_logo.svg/1200px-INDYCAR_logo.svg.png",
                "<@&1005573708590633063>"
            }},
            {"[IMSA]", new string[]
            {
                "665554362570899476",
                "https://upload.wikimedia.org/wikipedia/commons/thumb/9/94/IMSA_SportsCar_Championship_logo.svg/2560px-IMSA_SportsCar_Championship_logo.svg.png",
                "<@&1005574680486354964>"
            }},
            {"[NASCAR]", new string[]
            {
                "665554362570899476",
                "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cb/NASCAR_Cup_Series_logo.svg/1200px-NASCAR_Cup_Series_logo.svg.png",
                "<@&1005574994107044021>"
            }},
            {"[MotoGP]", new string[]
            {
                "665554362570899476",
                "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a0/Moto_Gp_logo.svg/1280px-Moto_Gp_logo.svg.png",
                "<@&1005573264619356290>"
            }},
            {"[NASA]", new string[]
            {
                "811641906685018172",
                "https://upload.wikimedia.org/wikipedia/commons/thumb/e/e5/NASA_logo.svg/2449px-NASA_logo.svg.png",
                "<@&1004652992416456754>"
            }},
            {"[SpaceX]", new string[]
            {
                "811641906685018172",
                "https://www.spacex.com/static/images/share.jpg",
                "<@&1004652992416456754>"
            }},
        };
    }

    class EventManager
    {
        private string eventsFile;
        private DataFrame events;
        private DataFrame compactEvents;

        public EventManager(string eventsFile)
        {
            this.eventsFile = eventsFile;
            this.events = DataFrame.LoadCsv(eventsFile, ',', false,
            new string[] 
            {
                "Category",
                "Title",
                "Description",
                "Date",
                "Channel",
                "Image",
                "Mention"
            },
            new Type[]
            {
                typeof(string),
                typeof(string),
                typeof(string),
                typeof(string),
                typeof(string),
                typeof(string),
                typeof(string)
            }
            );
            this.compactEvents = new DataFrame();
            UpdateCompactEvents();
        }

        public void Head(int lines)
        {
            Console.WriteLine(compactEvents.Head(lines));
        }

        public void Tail(int lines)
        {
            if (lines < 0 || lines > compactEvents.Rows.Count)
            {
                lines = Convert.ToInt32(compactEvents.Rows.Count);
            }
            Console.WriteLine(compactEvents.Tail(lines));
        }

        public void Insert(string category, string title, string description, string date, string channel, string image, string mention)
        {
            List<KeyValuePair<string, object>> row = new()
            {
                new KeyValuePair<string, object>("Category", category),
                new KeyValuePair<string, object>("Title", title),
                new KeyValuePair<string, object>("Description", description),
                new KeyValuePair<string, object>("Date", date),
                new KeyValuePair<string, object>("Channel", channel),
                new KeyValuePair<string, object>("Image", image),
                new KeyValuePair<string, object>("Mention", mention),
            };
            /*
            List<object> row = new ()
            {
                category,
                title,
                description,
                date,
                channel,
                "https://www.gluonspace.com/image.png",
                "<@23897918237>",
            };
            */
            events.Append(row, inPlace: true);
            events = events.OrderBy("Date");
            DataFrame.WriteCsv(events, eventsFile, ',', false);
            UpdateCompactEvents();
        }

        public void UpdateCompactEvents()
        {
            this.compactEvents = new DataFrame();
            compactEvents.Columns.Add(events.Columns[0]);
            compactEvents.Columns.Add(events.Columns[1]);
            compactEvents.Columns.Add(events.Columns[2]);
            compactEvents.Columns.Add(events.Columns[3]);
        }
    }
}

