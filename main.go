package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Results []struct {
		Gender string `json:"gender"`
		Name   struct {
			Title string `json:"title"`
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Location struct {
			Street struct {
				Number int    `json:"number"`
				Name   string `json:"name"`
			} `json:"street"`
			City        string `json:"city"`
			State       string `json:"state"`
			Country     string `json:"country"`
			Coordinates struct {
				Latitude  string `json:"latitude"`
				Longitude string `json:"longitude"`
			} `json:"coordinates"`
			Timezone struct {
				Offset      string `json:"offset"`
				Description string `json:"description"`
			} `json:"timezone"`
		} `json:"location"`
		Email string `json:"email"`
		Login struct {
			UUID     string `json:"uuid"`
			Username string `json:"username"`
			Password string `json:"password"`
			Salt     string `json:"salt"`
			Md5      string `json:"md5"`
			Sha1     string `json:"sha1"`
			Sha256   string `json:"sha256"`
		} `json:"login"`
		Dob struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"dob"`
		Registered struct {
			Date time.Time `json:"date"`
			Age  int       `json:"age"`
		} `json:"registered"`
		Phone string `json:"phone"`
		Cell  string `json:"cell"`
		ID    struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"id"`
		Picture struct {
			Large     string `json:"large"`
			Medium    string `json:"medium"`
			Thumbnail string `json:"thumbnail"`
		} `json:"picture"`
		Nat string `json:"nat"`
	} `json:"results"`
	Info struct {
		Seed    string `json:"seed"`
		Results int    `json:"results"`
		Page    int    `json:"page"`
		Version string `json:"version"`
	} `json:"info"`
}

func main() {
	if len(os.Args) != 3 {
		s := fmt.Sprintf("Welcome to passive v1.0.0\n\nOPTIONS:\n    -fn         Search with full-name\n    -ip         Search with ip address\n    -u          Search with username\n")
		fmt.Println(s)
		return
	}

	ip := flag.String("ip", "", "search by ip adress")
	username := flag.String("u", "", "search by fullname")
	fullname := flag.String("fn", "", "search by full name")

	flag.Parse()

	if len(*ip) != 0 {
		retreive_location(*ip)
		return
	}

	if len(*username) != 0 {
		checkUsername(*username)
		return
	}

	if len(*fullname) != 0 {
		retreive_profile(*fullname)
		return
	}
}

type Location struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	Org         string  `json:"org"`
	AS          string  `json:"as"`
	Query       string  `json:"query"`
}

type social_medias struct {
	twitter   bool
	instagram bool
	facebook  bool
	linkedin  bool
	skype     bool
}

func retreive_profile(fullname string) {
	if len(fullname) == 0 || len(strings.Split(fullname, " ")) != 2 {
		return
	}

	url := "https://randomuser.me/api/?results=1" // Number of users to retrieve

	// Create HTTP client
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	// Read and unmarshal the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var response User
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling response:", err)
		return
	}

	str := fmt.Sprintf("First Name: %s\nLast Name: %s\nAdress: %s %d %s\nNumber: %s", strings.Split(fullname, " ")[0], strings.Split(fullname, " ")[1], response.Results[0].Location.Street.Name, response.Results[0].Location.Street.Number, response.Results[0].Location.City, response.Results[0].Phone)

	isgood, result := writeFile(str)
	fmt.Println(str)
	if isgood {
		fmt.Println("Saved in", result)
	} else {
		fmt.Println("Error Writting file")
	}
}

func retreive_location(ip string) Location {
	if len(ip) == 0 || !isValidIPv4(ip) {
		fmt.Print("")

		return Location{}
	}

	if ip == "127.0.0.1" {

		// Créer une commande shell
		cmd := exec.Command("curl", "ipinfo.io/ip")

		// Exécuter la commande et récupérer la sortie
		output, err := cmd.Output()
		if err != nil {
			log.Fatalf("Erreur lors de l'exécution de la commande: %v", err)
		}

		ip = string(output)
	}

	// Check if it's a private IP address
	if ipNet := net.ParseIP(ip); ipNet != nil {
		if ipNet.IsPrivate() {
			fmt.Printf("IP Address: %s (Private IP)\n", ip)
			fmt.Println("This is a private IP address, typically used in local networks.")
			fmt.Println("City: N/A (Private Network)")
			fmt.Println("ISP: N/A (Private Network)")
			return Location{}
		}
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return Location{}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return Location{}
	}

	var ipInfo Location
	err = json.Unmarshal(body, &ipInfo)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return Location{}
	}

	s := fmt.Sprintf("ISP: %s\nCity Lat/Lon: (%.03f) / (%.03f)\nRegion: %s\n", ipInfo.ISP, ipInfo.Lat, ipInfo.Lon, ipInfo.Country)

	isgood, result := writeFile(s)

	fmt.Print(s)

	if isgood {
		fmt.Println("Saved in", result)
	} else {
		fmt.Println("Error Writting file")
	}

	return Location{}
}

func writeFile(content string) (bool, string) {
	baseFilename := "result.txt"

	// Find the next available filename
	filename := getNextAvailableFilename(baseFilename)

	// Create and write to the file
	file, err := os.Create(filename)
	if err != nil {
		return false, "Error creating file:"
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return false, "Error creating file:"
	}

	return true, filename

}

func getNextAvailableFilename(baseFilename string) string {
	filename := baseFilename
	counter := 1

	for {
		_, err := os.Stat(filename)
		if os.IsNotExist(err) {
			// File doesn't exist, we can use this filename
			return filename
		}

		// File exists, increment counter and try again
		parts := strings.Split(baseFilename, ".")
		nameWithoutExtension := parts[0]
		extension := "." + parts[1]

		filename = nameWithoutExtension + strconv.Itoa(counter) + extension
		counter++
	}
}

func isValidIPv4(ip string) bool {
	ipPattern := `^(\d{1,3}\.){3}\d{1,3}$`
	match, _ := regexp.MatchString(ipPattern, ip)
	if !match {
		return false
	}

	// Additional check for each octet
	octets := regexp.MustCompile(`\d+`).FindAllString(ip, -1)
	for _, octet := range octets {
		if i := len(octet); i > 1 && octet[0] == '0' {
			return false
		}
		if num := atoi(octet); num < 0 || num > 255 {
			return false
		}
	}
	return true
}

func atoi(s string) int {
	n := 0
	for _, ch := range s {
		ch -= '0'
		if ch > 9 {
			return -1
		}
		n = n*10 + int(ch)
	}
	return n
}

func checkUsername(username string) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Liste des URLs et indices de contenu à vérifier
	urls := map[string]struct {
		url     string
		keyword string
	}{
		"Github":   {fmt.Sprintf("https://github.com/%s", username), "Not Found"},
		"Snapchat": {fmt.Sprintf("https://www.snapchat.com/add/%s", username), "Page Not Found"},
		"Reddit":   {fmt.Sprintf("https://reddit.com/user/%s", username), "Sorry, nobody on Reddit"},
		"YouTube":  {fmt.Sprintf("https://www.youtube.com/%s", username), "This channel does not exist"},
		"TikTok":   {fmt.Sprintf("https://tiktok.com/@%s", username), "Couldn't find this account"},
	}

	count := 0
	str := ""

	for site, data := range urls {
		req, _ := http.NewRequest("GET", data.url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("%s: Erreur %v\n", site, err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		// Vérifier la présence de mots-clés d'erreur spécifiques dans la réponse
		if resp.StatusCode == 200 && !strings.Contains(string(body), data.keyword) {
			count++
			str += fmt.Sprintf("%s: yes\n", site)
			fmt.Print(fmt.Sprintf("%s: yes\n", site))
		} else {
			str += fmt.Sprintf("%s: no\n", site)
			fmt.Print(fmt.Sprintf("%s: no\n", site))
		}

		// Arrêter si on a trouvé le profil sur au moins 5 réseaux
		if count >= 5 {
			fmt.Println("Le nom d'utilisateur existe sur au moins 5 réseaux.")
			break
		}

		time.Sleep(1 * time.Second)
	}

	isgood, result := writeFile(str)

	if isgood {
		fmt.Println("Saved in", result)
	} else {
		fmt.Println("Error Writting file")
	}
}
