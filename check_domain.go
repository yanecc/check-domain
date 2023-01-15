package main

// 使用方法：ck -suffix app -path C:/domains.txt
import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
)

// type Config struct {
// 	ApiKey       string
// 	AccurateMode bool
// 	UseWhois     bool
// }

// type Check struct {
// 	Domain struct {
// 		Availability string `json:"domainAvailability"`
// 		Name         string `json:"domainName"`
// 	} `json:"DomainInfo"`
// }

// type Balance struct {
// 	Data []struct {
// 		ProductID int `json:"product_id"`
// 		Product   struct {
// 			ID   int    `json:"id"`
// 			Name string `json:"name"`
// 		} `json:"product"`
// 		Credits int `json:"credits"`
// 	} `json:"data"`
// }

var (
	config     Config
	suffix     string
	path       string
	fileName   = "domains.txt"
	configFile = "config.toml"
	okFile, _  = os.Create("ok.txt")
	okWriter   = bufio.NewWriter(okFile)
)

func init() {
	flag.StringVar(&path, "path", "", "The path of the domain name prefix file")
	flag.StringVar(&suffix, "suffix", "", "The domain name suffix")
}

func main() {
	defer func(okFile *os.File) {
		err := okFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(okFile)
	defer func(okWriter *bufio.Writer) {
		err := okWriter.Flush()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(okWriter)
	flag.Parse()

	// read config from config.toml
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		fmt.Println("config.toml not found, creating...")
		createConfigFile()
		os.Exit(0)
	}

	// read domain name prefix from file
	var data []byte
	var err error
	if path == "" {
		data, err = os.ReadFile(fileName)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	lines := strings.Split(string(data), "\n")
	var prefixes []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			prefixes = append(prefixes, line)
		}
	}
	total := len(prefixes)

	// check domain name availability
	for i, prefix := range prefixes {
		var domain string
		if suffix == "" {
			domain = prefix
		} else {
			domain = prefix + "." + suffix
		}
		fmt.Printf("Checking: %s (%d/%d)", domain, i+1, total)
		url := fmt.Sprintf("https://domain-availability.whoisxmlapi.com/api/v1?apiKey=%s&domainName=%s&credits=%s&mode=%s", config.ApiKey, domain, getCredits(), getMode())
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(resp.Body)

		var result Check
		err = json.NewDecoder(resp.Body).Decode(&result)
		if err != nil {
			return
		}

		if result.Domain.Availability == "AVAILABLE" {
			color.Green("  ✓")
			_, err := okWriter.WriteString(result.Domain.Name + "\n")
			if err != nil {
				return
			}
		} else {
			fmt.Println()
		}
	}
	fmt.Println("\nFinished writing available domains to ok.txt")
	checkBalance()
}

// func createConfigFile() {
// 	config := Config{
// 		ApiKey:       "apiKey",
// 		AccurateMode: true,
// 		UseWhois:     false,
// 	}

// 	configFile, err := os.Create(configFile)
// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// 	defer configFile.Close()

// 	encoder := toml.NewEncoder(configFile)
// 	if err := encoder.Encode(config); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

//		fmt.Println("Please register and get api key from https://user.whoisxmlapi.com/products, then edit config.toml.")
//	}
// 	fmt.Println("Please register and get api key from https://user.whoisxmlapi.com/products, then edit config.toml.")
// }

func createConfigFile() {
	configFile, err := os.Create(configFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(configFile)

	_, err = fmt.Fprintln(configFile, "# api获取地址：https://user.whoisxmlapi.com/products")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = fmt.Fprintln(configFile, "apiKey = \"apiKey\"")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = fmt.Fprintln(configFile, "\naccurateMode = true")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = fmt.Fprintln(configFile, "\n# 免费账户有500次Whois查询额度及100次/月域名可用性检测额度")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = fmt.Fprintln(configFile, "useWhois = false")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Please register and get api key from https://user.whoisxmlapi.com/products, then edit config.toml.")
}

func getMode() string {
	if config.AccurateMode {
		return "DNS_AND_WHOIS"
	}
	return "DNS_ONLY"
}

func getCredits() string {
	if config.UseWhois {
		return "WHOIS"
	}
	return "DA"
}

func checkBalance() {
	url := fmt.Sprintf("https://user.whoisxmlapi.com/user-service/account-balance?apiKey=%s", config.ApiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(resp.Body)

	var balance Balance
	err = json.NewDecoder(resp.Body).Decode(&balance)
	if err != nil {
		return
	}

	var whois, da int
	for _, v := range balance.Data {
		if v.Product.Name == "WHOIS API" {
			whois = v.Credits
		} else if v.Product.Name == "Domain Availability API" {
			da = v.Credits
		}
	}
	fmt.Printf("您的 Whois 查询次数余额：%d\n", whois)
	fmt.Printf("域名可用性检测次数余额：%d/100（每月重置）\n", da)
}
