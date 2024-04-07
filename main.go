package main

import (
	Auth "CLI/Auth"
	HugginFace "CLI/HugAI"
	Searx "CLI/SearXNG"
	Movie_ "CLI/TMDB"
	blackbox "CLI/blackbox"
	MerlinAI "CLI/merlin_cli"
	"CLI/s7cli/commands"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func openUrlInDefaultBrowser(url string) {
	var cmdName string

	switch {
	case isUnix():
		cmdName = "xdg-open"
	case isMacOs():
		cmdName = "open"
	default:
		cmdName = "start"
	}

	cmdArgs := []string{"cmd.exe", "/C", cmdName, url}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)

	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Error occurred starting the command '%s'. Error:\n%s", strings.Join(cmdArgs, " "), err.Error())
		fmt.Println("\nOutput:", string(output[:]))
	} else {
		fmt.Printf("Successfully started command '%s'\n", strings.Join(cmdArgs, " "))
	}
}

func isUnix() bool {
	return sysInfo.GOOS == "linux" || sysInfo.GOOS == "darwin"
}

func isMacOs() bool {
	return sysInfo.GOOS == "darwin"
}

type SystemInfo struct {
	GOARCH string
	GOOS   string
}

var sysInfo = SystemInfo{}

func init() {
	sysInfo.GOARCH = os.Getenv("GOARCH")
	sysInfo.GOOS = os.Getenv("GOOS")
}

func merlin(args []string, this commands.Command) error {

	settings_file, err := os.Open("settings.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return nil
	}

	data, err := ioutil.ReadAll(settings_file)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return nil
	}

	type Data struct {
		Key1 string `json:"merlin_auth_token"`
		Key2 string `json:"huggingface_cookie"`
	}

	var result Data
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}
	defer settings_file.Close()

	authToken := result.Key1
	chatID := "43ac5495-e1e1-4a68-9115-" + this.Name
	m := MerlinAI.NewMerlin(authToken, chatID)
	message := strings.Join(args[1:], " ")

	responseBody, err := m.Chat(message)
	if err != nil {
		log.Fatal(err)
	}

	err = m.StreamContent(responseBody)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
func set_window_title(title string) any {
	cmd := exec.Command("cmd", "/C", "title", title)
	cmd.Run()
	return nil
}

func promptString(fcname string) string {
	now := time.Now()

	currentTime := now.Format("15:04:05")

	return fmt.Sprintf("[%s][%s]@[CLI]~# ", currentTime, fcname)
}

func Banner() {
	Banner := commands.Reset + `
	███████╗███████╗     █████╗ ██╗     ██████╗██╗     ██╗
	██╔════╝╚════██║    ██╔══██╗██║    ██╔════╝██║     ██║
	███████╗    ██╔╝    ███████║██║    ██║     ██║     ██║
	╚════██║   ██╔╝     ██╔══██║██║    ██║     ██║     ██║
	███████║   ██║      ██║  ██║██║    ╚██████╗███████╗██║
	╚══════╝   ╚═╝      ╚═╝  ╚═╝╚═╝     ╚═════╝╚══════╝╚═╝
	`
	ColoredBanner := strings.ReplaceAll(Banner, "█", commands.BoldPurple+"█"+commands.Reset)
	fmt.Println(ColoredBanner)
}

var auth_status bool = false

func main() {
	osargs := os.Args
	Banner()

	settings_file, err := os.Open("settings.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}

	data, err := ioutil.ReadAll(settings_file)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	type Data struct {
		Key1     string `json:"merlin_auth_token"`
		Key2     string `json:"huggingface_cookie"`
		Key3     string `json:"blackbox_cookie"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var result Data
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	defer settings_file.Close()
	set_window_title("CLI")
	DefaultHandler := commands.DefaultHandler
	type Command = commands.Command
	type Arg = commands.Arg
	DefaultHandler.AddCommand(Command{
		Name:        "exit",
		Description: "Exits this application.",
		Args:        []Arg{},
		Exec: func(args []string, command Command) error {
			if isUnix() {
				fmt.Println("Exiting... " + osargs[1])

				output, err := exec.Command("/bin/getppid.sh", osargs[1]).Output()
				if err != nil {
					log.Fatalf("Failed to execute command: %v", err)
				}

				result := strings.Split(string(output[:]), "\n")

				fmt.Println("Result:", result[0])
				process, err := os.StartProcess("/bin/killwtsk.sh", []string{"/bin/killwtsk.sh", strings.Split(result[0], "?")[0]}, &os.ProcAttr{Env: os.Environ()})
				if err == nil {

					err = process.Release()
					if err != nil {
						fmt.Println(err.Error())
					}

				} else {
					fmt.Println(err.Error())
				}
			}
			if !isUnix() {
				os.Exit(0)
			}
			return nil

		},
	})

	authService := Auth.AuthService{}
	for {
		if result.Username != "" || result.Password != "" {
			break
		}
		DefaultHandler.SetPrompt("> ")
		fmt.Println("Enter your username:")
		username := DefaultHandler.GetInput()
		fmt.Println("Enter your password:")
		password := DefaultHandler.GetInput()
		fmt.Println("Logging in...")
		result.Username = username
		result.Password = password
		file, err := os.OpenFile("settings.json", os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println("Error opening JSON file:", err)
			return
		}
		datax, err := json.MarshalIndent(&result, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}
		defer file.Close()
		err = ioutil.WriteFile("settings.json", datax, 0644)
		if err != nil {
			fmt.Println("Error writing to JSON file:", err)
			return
		}
	}
	response, err := authService.PerformLogin(result.Username, result.Password)
	if err != nil {
		fmt.Printf("An error occurred: %+v\n", err)
		return
	}
	switch response.StatusCode {
	case 200:
		authenticated, _ := response.Data.(map[string]interface{})["authenticated"].(bool)
		fmt.Println("Authenticated:", authenticated)
		username, _ := response.Data.(map[string]interface{})["username"].(string)
		fmt.Println("Username:", username)
		role, _ := response.Data.(map[string]interface{})["role"].(string)
		fmt.Println("Role:", role)
	default:
		set_window_title("CLI > Unauthenticated")
		DefaultHandler.AddCommand(commands.Command{
			Name:        "login",
			Description: "Login to the application.",
			Exec: func(input []string, this commands.Command) error {
				fmt.Println("Enter your username:")
				username := DefaultHandler.GetInput()
				fmt.Println("Enter your password:")
				password := DefaultHandler.GetInput()
				response, err := authService.PerformLogin(username, password)
				if err != nil {
					fmt.Printf("An error occurred: %+v\n", err)
					return nil
				}
				if response.StatusCode != 200 {
					fmt.Println("Authentication failed with status code", response.StatusCode)
					return nil
				}
				if response.StatusCode == 200 {
					auth_status = true
					authenticated, _ := response.Data.(map[string]interface{})["authenticated"].(bool)
					fmt.Println("Authenticated:", authenticated)
					username, _ := response.Data.(map[string]interface{})["username"].(string)
					fmt.Println("Username:", username)
					role, _ := response.Data.(map[string]interface{})["role"].(string)
					fmt.Println("Role:", role)
					settings_file, err := os.Open("settings.json")
					if err != nil {
						fmt.Println("Error opening JSON file:", err)
						return nil
					}
					data, err := ioutil.ReadAll(settings_file)
					if err != nil {
						fmt.Println("Error reading JSON file:", err)
						return nil
					}
					type Data_ struct {
						Key1     string `json:"merlin_auth_token"`
						Key2     string `json:"huggingface_cookie"`
						Key3     string `json:"blackbox_cookie"`
						Username string `json:"username"`
						Password string `json:"password"`
					}
					var result_ Data_
					if err := json.Unmarshal(data, &result_); err != nil {
						fmt.Println("Error unmarshalling JSON:", err)
						return nil
					}
					defer settings_file.Close()

					result_.Username = username
					result_.Password = password
					data, err = json.MarshalIndent(&result_, "", "  ")
					if err != nil {
						fmt.Println("Error marshalling JSON:", err)
						return nil
					}
					err = ioutil.WriteFile("settings.json", data, 0644)
					if err != nil {
						fmt.Println("Error writing JSON file:", err)
						return nil
					}
					fmt.Println("Settings file updated successfully")
				}
				return nil
			},
		})
		set_window_title("CLI > Unauthenticated Menu")
		DefaultHandler.SetPrompt(promptString("NoAuth/"))
		for {

			if !auth_status {
				fmt.Println("Authentication failed")

			}
			if auth_status {
				break
			}
			input_x := DefaultHandler.GetInput()
			if input_x == "exit" {
				return
			}
			if input_x == "clear" {
				DefaultHandler.Handle(input_x)
				continue
			}
			DefaultHandler.Handle(input_x)
		}
	}
	set_window_title("CLI > User: " + result.Username)
	DefaultHandler.SetPrompt(promptString("base/"))

	DefaultHandler.AddCommand(Command{
		Name:        "merlin",
		Description: "Merlin AI (GPT 3)",
		Args:        []Arg{},
		Exec: func(input []string, this Command) error {
			for {
				DefaultHandlerx := commands.DefaultHandler
				DefaultHandlerx.SetPrompt(promptString("base/Merlin"))
				DefaultHandlerx.AddCommand(commands.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				message := DefaultHandlerx.GetInput()

				if message == "exit" {
					break

				}
				if message == "" {
					continue
				}
				if message == "clear" {
					DefaultHandlerx.Handle(message)
					continue
				}
				x := strings.Split(message, " ")
				merlin(x, this)
			}

			return nil
		},
	})
	DefaultHandler.AddCommand(Command{
		Name:        "hug",
		Description: "Hugging AI (?)",
		Args:        []Arg{},
		Exec: func(input []string, this Command) error {
			client := HugginFace.NewHug()
			ChatId := "6608a05392dfb775db102588"
			cookie := result.Key2
			for {
				DefaultHandlerx2 := commands.DefaultHandler
				DefaultHandlerx2.SetPrompt(promptString("base/Hug"))
				DefaultHandlerx2.AddCommand(commands.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				DefaultHandlerx2.AddCommand(commands.Command{
					Name:        "change model",
					Description: "Change the current AI Model",
				})
				message := DefaultHandlerx2.GetInput()

				if message == "exit" {
					break
				}
				if message == "" {
					continue
				}
				if message == "clear" {
					DefaultHandlerx2.Handle(message)
					continue
				}
				if message == "change model" {
					fmt.Println("Available Models:")
					fmt.Println("Use the TAB Button on your keyboard to cycle through the list of models")
					Df := commands.DefaultHandler
					Df.SetPrompt("> ")
					Df.AddCommand(commands.Command{
						Name:        "google/gemma-7b-it",
						Description: "Google AI",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})

					Df.AddCommand(commands.Command{
						Name:        "mistralai/Mixtral-8x7B-Instruct-v0.1",
						Description: "Mixtral Chat AI v0.1",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})

					Df.AddCommand(commands.Command{
						Name:        "mistralai/Mistral-7B-Instruct-v0.2",
						Description: "Mixtral Chat AI v0.2",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands.Command{
						Name:        "meta-llama/Llama-2-70b-chat-hf",
						Description: "Facebook (Meta) Llama AI",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands.Command{
						Name:        "NousResearch/Nous-Hermes-2-Mixtral-8x7B-DPO",
						Description: "NousResearch x Mixtral-8x7B",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands.Command{
						Name:        "codellama/CodeLlama-70b-Instruct-hf",
						Description: "CodeLlama (Programming Assistant AI)",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands.Command{
						Name:        "openchat/openchat-3.5-0106",
						Description: "OpenChat 3.5 (GPT 3.5 Turbo)",
						Exec: func(input []string, this commands.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					for {
						input_ := Df.GetInput()

						if input_ == "" {
							continue
						} else if input_ == "clear" {
							Df.Handle(input_)
						} else {
							fmt.Printf("Model has been set to [%s]\r\n", input_)
							Df.Handle(
								input_,
							)
							break
						}

					}
					continue
				}

				Id_ := client.GetMsgUID(ChatId, cookie)
				err := client.SendMessage(message, ChatId, Id_, cookie)
				if err != nil {
					log.Fatal(err)
				}
			}
			return nil
		},
	})

	DefaultHandler.AddCommand(Command{
		Name:        "blackbox",
		Description: "BlackBox Programming AI Chat",
		Args:        []Arg{},
		Exec: func(input []string, this commands.Command) error {
			for {
				DefaultHandlerx3 := commands.DefaultHandler
				DefaultHandlerx3.SetPrompt(promptString("base/BlackBox"))
				DefaultHandlerx3.AddCommand(commands.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				message := DefaultHandlerx3.GetInput()
				if message == "exit" {
					break
				}
				if message == "" {
					continue
				}
				if message == "clear" {
					DefaultHandlerx3.Handle(message)
					continue
				}

				BlackBox_ := blackbox.NewBlackboxClient()
				BlackBox_.SendMessage(message)
			}
			return nil
		},
	})

	type Item struct {
		Href  string `json:"href"`
		Desc  string `json:"desc"`
		Title string `json:"title"`
	}
	DefaultHandler.AddCommand(Command{
		Name:        "searx",
		Description: "Use Searx Search Engine",
		Args:        []Arg{},
		Exec: func(input []string, this commands.Command) error {
			for {
				DefaultHandlerx4 := commands.DefaultHandler
				DefaultHandlerx4.SetPrompt(promptString("base/Searx"))
				DefaultHandlerx4.AddCommand(commands.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				DefaultHandlerx4.AddCommand(commands.Command{
					Name:        "search",
					Description: "Search for something on Searx",
				})
				message := DefaultHandlerx4.GetInput()
				if message == "exit" {
					break
				}
				if message == "" {
					continue
				}
				if message == "clear" {
					DefaultHandlerx4.Handle(message)
					continue
				}

				if message == "search" {
					Searx_ := Searx.NewSearchEngine()
					DefaultHandlerx4.SetPrompt(promptString("base/Searx/Search"))
					DefaultHandlerx4.AddCommand(commands.Command{
						Name:        "exit",
						Description: "Exit back to the main cli",
					})
					message := DefaultHandlerx4.GetInput()
					if message == "" {
						fmt.Println("Please enter a search query")
					}
					if message != "" {
						results_ := Searx_.Run(message)

						var results []Item
						err = json.Unmarshal([]byte(results_), &results)
						if err != nil {
							log.Fatal(err)
						}
						x := 0
						DefaultHandlerx4.SetPrompt(promptString("base/Searx/Search/Results"))
						for i := range results {
							results[i].Title = strings.ReplaceAll(results[i].Title, "\n", "")
							if len(results[i].Title) > 100 {
								results[i].Title = results[i].Title[0:100]
							}
							x = x + 1
							DefaultHandlerx4.AddCommand(commands.Command{
								Name:        `>` + strconv.Itoa(x),
								Description: results[i].Title,
								Args:        []Arg{},
								Exec: func(input []string, this commands.Command) error {
									DefaultHandlerx4.Handle("clear")
									fmt.Println("==============================================================")
									fmt.Println(results[i].Href)
									fmt.Println(results[i].Desc)
									fmt.Println(results[i].Title)
									fmt.Println("==============================================================")
									return nil
								},
							})

						}
						result_picker := DefaultHandlerx4.GetInput()
						if result_picker == "exit" {
							break
						}
						if result_picker == "clear" {
							DefaultHandlerx4.Handle(result_picker)
							continue
						}
						if result_picker == "" {
							continue
						}
						DefaultHandlerx4.Handle(result_picker)
					}

				}

			}
			return nil
		},
	})
	DefaultHandler.AddCommand(Command{
		Name:        "movie",
		Description: "Search for a movie",
		Args:        []Arg{},
		Exec: func(input []string, this commands.Command) error {
			for {
				DefaultHandlerx5 := commands.DefaultHandler
				DefaultHandlerx5.SetPrompt(promptString("base/Movie"))
				DefaultHandlerx5.AddCommand(commands.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				DefaultHandlerx5.AddCommand(commands.Command{
					Name:        "search",
					Description: "Search for a movie",
				})
				message := DefaultHandlerx5.GetInput()
				if message == "exit" {
					break
				}
				if message == "" {
					continue
				}
				if message == "clear" {
					DefaultHandlerx5.Handle(message)
					continue
				}
				if message == "search" {
					tmdb_ := Movie_.Init("71e68428e0a8d7f642158c4cc4c74f4c")
					DefaultHandlerx5.SetPrompt(promptString("base/Movie/Search"))
					message := DefaultHandlerx5.GetInput()
					resp, err := tmdb_.MovieData(message)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Found " + strconv.Itoa(len(resp.Results)) + " results")
					if len(resp.Results) == 0 {
						fmt.Println("No results found")
						continue
					}
					fmt.Println("Press to view the results")
					for i := range resp.Results {

						DefaultHandlerx5.AddCommand(commands.Command{
							Name:        ">" + strconv.Itoa(i),
							Description: resp.Results[i].Title,
							Exec: func(input []string, this commands.Command) error {
								DefaultHandlerx5.Handle("clear")
								fmt.Println("> " + resp.Results[i].Title)
								fmt.Println("==============================================================")
								fmt.Println("https://vidsrc.xyz/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id))
								fmt.Println("==============================================================")
								fmt.Println("==============================================================")
								fmt.Println("https://vidsrc.net/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id))
								fmt.Println("==============================================================")
								fmt.Println("==============================================================")
								fmt.Println("https://vidsrc.in/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id))
								fmt.Println("==============================================================")
								internal_handler := commands.DefaultHandler
								internal_handler.SetPrompt(promptString("base/Movie/Search/Results"))
								internal_handler.AddCommand(commands.Command{
									Name:        "exit",
									Description: "Exit back to the main cli",
								})
								list_ := []string{"https://vidsrc.xyz/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id), "https://vidsrc.in/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id), "https://vidsrc.net/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id)}

								for i := range list_ {
									internal_handler.AddCommand(commands.Command{
										Name:        ">" + strconv.Itoa(i),
										Description: list_[i],
										Exec: func(input []string, this commands.Command) error {
											DefaultHandlerx5.Handle("clear")
											fmt.Println("> " + list_[i])
											fmt.Println("==============================================================")
											openUrlInDefaultBrowser(list_[i])

											return nil
										},
									})

								}
								fmt.Println("Press tab to get the list of results")
								internal_handler.SetPrompt(promptString("base/Movie/Search/Results"))
								message := internal_handler.GetInput()
								internal_handler.Handle(message)
								return nil
							},
						})
					}
					result_picker := DefaultHandlerx5.GetInput()
					if result_picker == "exit" {
						break
					}
					if result_picker == "clear" {
						DefaultHandlerx5.Handle(result_picker)
						continue
					}
					if result_picker == "" {
						continue
					}
					DefaultHandlerx5.Handle(result_picker)

				}
			}
			return nil
		},
	})

	for {
		DefaultHandler.SetPrompt(promptString("base/"))
		input_data := DefaultHandler.GetInput()

		DefaultHandler.Handle(
			input_data,
		)
		if input_data == "clear" {
			Banner()
		}
	}
}
