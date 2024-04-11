package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	commands2 "main/pkg/commands"
	"main/pkg/misc"
	BlackBox "main/pkg/utils/blackbox"
	HugginFace "main/pkg/utils/huggingface"
	Searx "main/pkg/utils/searx"
	Movie_ "main/pkg/utils/tmdb"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

type MC struct {
}

type Searx_Struct struct {
	Href  string `json:"href"`
	Desc  string `json:"desc"`
	Title string `json:"title"`
}

var (
	f_            = misc.Funcs{}
	settings, err = f_.LoadSettings()
)

func (m *MC) Run(h commands2.Handler) {
	for {
		h.SetPrompt("> ")
		handler_input := h.GetInput()
		h.Handle(handler_input)
	}

}
func (m *MC) GetInput() string {
	DefaultHandler := commands2.DefaultHandler
	DefaultHandler.SetPrompt("~# ")
	return DefaultHandler.GetInput()

}
func (m *MC) Init(h commands2.Handler) commands2.Handler {

	type Command = commands2.Command
	type Arg = commands2.Arg
	// Clear command
	h.AddCommand(commands2.Command{
		Name:        "clear",
		Description: "Clears the console.",
		Args:        []Arg{},
		Exec: func(input []string, this Command) error {
			os_switch := make(map[string]func()) //Initialize it
			os_switch["linux"] = func() {
				cmd := exec.Command("clear") //Linux example, its tested
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
			os_switch["windows"] = func() {
				cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
				cmd.Stdout = os.Stdout
				cmd.Run()
			}

			value, ok := os_switch[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
			if ok {                              //if we defined a clear func for that platform:
				value() //we execute it
			} else { //unsupported platform
				fmt.Println("Failed; Your terminal isn't ANSI! :(")
			}
			f_.Banner()
			return nil
		},
	})
	h.AddCommand(Command{
		Name:        "merlin",
		Description: "Merlin AI (GPT 3)",
		Args:        []Arg{},
		Exec: func(input []string, this Command) error {
			for {
				DefaultHandlerx := commands2.DefaultHandler
				DefaultHandlerx.SetPrompt("Merlin > ")
				DefaultHandlerx.AddCommand(commands2.Command{
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
				f_.MerlinAI_(x, this)
			}

			return nil
		},
	})
	h.AddCommand(Command{
		Name:        "hug",
		Description: "Hugging AI (?)",
		Args:        []Arg{},
		Exec: func(input []string, this Command) error {
			client := HugginFace.NewHug()
			ChatId := "6608a05392dfb775db102588"
			cookie := settings.HugginFaceCookie
			for {
				DefaultHandlerx2 := commands2.DefaultHandler
				DefaultHandlerx2.SetPrompt("Hugging Face > ")
				DefaultHandlerx2.AddCommand(commands2.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				DefaultHandlerx2.AddCommand(commands2.Command{
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
					Df := commands2.DefaultHandler
					Df.SetPrompt("> ")
					Df.AddCommand(commands2.Command{
						Name:        "google/gemma-7b-it",
						Description: "Google AI",
						Exec: func(input []string, this commands2.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})

					Df.AddCommand(commands2.Command{
						Name:        "mistralai/Mixtral-8x7B-Instruct-v0.1",
						Description: "Mixtral Chat AI v0.1",
						Exec: func(input []string, this commands2.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})

					Df.AddCommand(commands2.Command{
						Name:        "mistralai/Mistral-7B-Instruct-v0.2",
						Description: "Mixtral Chat AI v0.2",
						Exec: func(input []string, this commands2.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands2.Command{
						Name:        "meta-llama/Llama-2-70b-chat-hf",
						Description: "Facebook (Meta) Llama AI",
						Exec: func(input []string, this commands2.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands2.Command{
						Name:        "NousResearch/Nous-Hermes-2-Mixtral-8x7B-DPO",
						Description: "NousResearch x Mixtral-8x7B",
						Exec: func(input []string, this commands2.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands2.Command{
						Name:        "codellama/CodeLlama-70b-Instruct-hf",
						Description: "CodeLlama (Programming Assistant AI)",
						Exec: func(input []string, this commands2.Command) error {
							model_ := this.Name
							ChatId = client.ChangeModel(model_, cookie)
							return nil
						},
					})
					Df.AddCommand(commands2.Command{
						Name:        "openchat/openchat-3.5-0106",
						Description: "OpenChat 3.5 (GPT 3.5 Turbo)",
						Exec: func(input []string, this commands2.Command) error {
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
	h.AddCommand(Command{
		Name:        "blackbox",
		Description: "BlackBox Programming AI Chat",
		Args:        []Arg{},
		Exec: func(input []string, this commands2.Command) error {
			for {
				DefaultHandlerx3 := commands2.DefaultHandler
				DefaultHandlerx3.SetPrompt("BlackBox > ")
				DefaultHandlerx3.AddCommand(commands2.Command{
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

				BlackBox_ := BlackBox.NewBlackboxClient()
				BlackBox_.SendMessage(message)
			}
			return nil
		},
	})
	h.AddCommand(Command{
		Name:        "searx",
		Description: "Use Searx Search Engine",
		Args:        []Arg{},
		Exec: func(input []string, this commands2.Command) error {
			for {
				DefaultHandlerx4 := commands2.DefaultHandler
				DefaultHandlerx4.SetPrompt("Searx > ")
				DefaultHandlerx4.AddCommand(commands2.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				DefaultHandlerx4.AddCommand(commands2.Command{
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
					DefaultHandlerx4.SetPrompt("Searx Search > ")
					DefaultHandlerx4.AddCommand(commands2.Command{
						Name:        "exit",
						Description: "Exit back to the main cli",
					})
					message := DefaultHandlerx4.GetInput()
					if message == "" {
						fmt.Println("Please enter a search query")
					}
					if message != "" {
						results_ := Searx_.Run(message)

						var results []Searx_Struct
						err = json.Unmarshal([]byte(results_), &results)
						if err != nil {
							log.Fatal(err)
						}
						x := 0
						DefaultHandlerx4.SetPrompt("Searx Results > ")
						for i := range results {
							results[i].Title = strings.ReplaceAll(results[i].Title, "\n", "")
							if len(results[i].Title) > 100 {
								results[i].Title = results[i].Title[0:100]
							}
							x = x + 1
							DefaultHandlerx4.AddCommand(commands2.Command{
								Name:        `>` + strconv.Itoa(x),
								Description: results[i].Title,
								Args:        []Arg{},
								Exec: func(input []string, this commands2.Command) error {
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
	h.AddCommand(Command{
		Name:        "movie",
		Description: "Search for a movie",
		Args:        []Arg{},
		Exec: func(input []string, this commands2.Command) error {
			for {
				DefaultHandlerx5 := commands2.DefaultHandler
				DefaultHandlerx5.SetPrompt("Movie > ")
				DefaultHandlerx5.AddCommand(commands2.Command{
					Name:        "exit",
					Description: "Exit back to the main cli",
				})
				DefaultHandlerx5.AddCommand(commands2.Command{
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
					DefaultHandlerx5.SetPrompt("Movie Search > ")
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

						DefaultHandlerx5.AddCommand(commands2.Command{
							Name:        ">" + strconv.Itoa(i),
							Description: resp.Results[i].Title,
							Exec: func(input []string, this commands2.Command) error {
								h.Handle("clear")

								title_ := resp.Results[i].Title
								list_ := []string{"https://vidsrc.xyz/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id), "https://vidsrc.in/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id), "https://vidsrc.net/embed/movie?tmdb=" + strconv.Itoa(resp.Results[i].Id)}

								poster := "https://image.tmdb.org/t/p/original" + resp.Results[i].Poster_path
								image_format := poster[len(poster)-3:]
								var source_url string
								var yN bool
								// var base *huh.Theme = huh.ThemeBase()
								var dracula *huh.Theme = huh.ThemeCharm()
								// var base16 *huh.Theme = huh.ThemeBase16()
								// var charm *huh.Theme = huh.ThemeCharm()
								// var catppuccin *huh.Theme = huh.ThemeCatppuccin()
								form := huh.NewForm(
									huh.NewGroup(
										huh.NewSelect[string]().
											Options(huh.NewOptions(list_...)...).
											Title(title_).Value(&source_url),
										huh.NewConfirm().
											Title("Are you sure? ").
											Description("Please confirm. ").
											Affirmative("Yes!").
											Negative("No.").
											Inline(true).
											Value(&yN),
									),
								).WithAccessible(true).WithTheme(dracula)

								err := form.Run()
								if err != nil {
									log.Fatal(err)
								}
								// fmt.Println("> " + source_url)

								if image_format == "jpg" || image_format == "png" {
									fmt.Println("Poster: " + poster)
								} else {
									fmt.Println("No poster available")
								}
								//imageb64 := f_.ImageURLToBase64(poster)
								//img(image_format, title_, imageb64)
								if yN {
									f_.OpenUrl(source_url)
								}
								internal_handler := commands2.DefaultHandler
								internal_handler.SetPrompt("Back? > ")
								internal_handler.AddCommand(commands2.Command{
									Name:        "exit",
									Description: "Exit back to the main cli",
								})

								message := internal_handler.GetInput()
								if message == "exit" {
									return nil
								}
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
	h.AddCommand(Command{
		Name:        "settings",
		Description: "Configure application settings",
		Args:        []Arg{},
		Exec: func(input []string, this commands2.Command) error {
			f_.SettingsPage()
			f_.Banner()
			return nil
		},
	})

	return h
}
