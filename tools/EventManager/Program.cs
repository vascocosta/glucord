namespace EventManager
{
    class Program
    {   
        static void Main(string[] args)
        {
            if (args.Length != 1)
            {
                ShowUsage();
                Environment.Exit(1);
            }
            EventManager? eventManager = null;
            try
            {
                eventManager = new EventManager(args[0]);
            }
            catch (FileNotFoundException)
            {
                Console.WriteLine("Events file not found.");
                Environment.Exit(1);
            }
            while (true)
            {
                Console.Write("Command> ");
                string? input = Console.ReadLine();
                if (input == null)
                {
                    continue;
                }
                string[] splitInput = input.Split();
                switch (splitInput[0])
                {
                    case "head":
                        if (splitInput.Length == 2)
                        {
                            try {
                                int lines = Convert.ToInt32(splitInput[1]);
                                eventManager.Head(lines);
                            }
                            catch (FormatException)
                            {
                                Console.WriteLine("Wrong lines format.");
                            }
                        }
                        else {
                            Console.WriteLine("Wrong syntax.");
                        }
                        break;
                    case "tail":
                        if (splitInput.Length == 2)
                        {
                            try {
                                int lines = Convert.ToInt32(splitInput[1]);
                                eventManager.Tail(lines);
                            }
                            catch (FormatException)
                            {
                                Console.WriteLine("Wrong lines format.");
                            }
                        }
                        else {
                            Console.WriteLine("Wrong syntax.");
                        }
                        break;
                    case "insert":
                    case "new":
                        Console.Write("Category> ");
                        string? category = Console.ReadLine();
                        Console.Write("Title> ");
                        string? title = Console.ReadLine();
                        Console.Write("Description> ");
                        string? description = Console.ReadLine();
                        Console.Write("Date> ");
                        string? date = Console.ReadLine();
                        if (category != null && title != null && description != null && date != null)
                        {
                            Console.WriteLine($"Category: {category}\nTitle: {title}\nDescription: {description}\nDate: {date}");
                            Console.Write("Confirm> ");
                            input = Console.ReadLine();
                            Console.Write("Inserting new event... ");
                            try
                            {
                                eventManager.Insert(category, title, description, date, Resolver.categories[category][0], Resolver.categories[category][1], Resolver.categories[category][2]);
                                Console.WriteLine("OK");
                            }
                            catch (KeyNotFoundException)
                            {
                                Console.WriteLine("Fail: Invalid category.");
                            }
                        }
                        else
                        {
                            Console.WriteLine("Couldn't read input correctly.");
                        }
                        break;
                    case "exit":
                    case "quit":
                         Environment.Exit(0);
                        break;
                    default:
                        Console.WriteLine("Available commands:\nhead <lines>\ntail <lines>\ninsert/new");
                        break;
                }
            }
        }

        static void ShowUsage()
        {
            Console.WriteLine("Usage: EventManager <events_file>");

        }
    }
}
