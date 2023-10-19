package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mattn/go-sixel"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("* xfb_pair_generate_text2stickers *")
	fmt.Print("prompt: ")

	var prompt string
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	prompt = scanner.Text()

	generateImg(prompt)
}

func generateImg(prompt string) {
	client := http.Client{}

	payload := "method=post&pretty=false&format=json&server_timestamps=true&locale=en_US&fb_api_req_friendly_name=IGGenerateText2Stickers&client_doc_id=272961684716926046129519059692&enable_canonical_naming=true&enable_canonical_variable_overrides=true&enable_canonical_naming_ambiguous_type_prefixing=true&variables=%7B%22is_pando%22%3Atrue%2C%22bypass_cache%22%3Atrue%2C%22prompt%22%3A%22" + strings.ReplaceAll(prompt, " ", "+") + "%22%2C%22caller%22%3A%22ig_stories_ai_stickers%22%2C%22media_type%22%3A%22image%2Fpng%22%2C%22scaling_factor%22%3A1%7D"

	req, _ := http.NewRequest("POST", "https://i.instagram.com/graphql_www", strings.NewReader(payload))
	req.Header.Set("User-Agent", "Instagram 302.1.0.36.111 Android (30/11; 180dpi; 1011x1495; Waydroid/waydroid; WayDroid x86_64 Device; waydroid_x86_64; unknown; en_US; 520702298)")
	req.Header.Set("X-Tigon-Is-Retry", "False")
	req.Header.Set("X-Fb-Rmd", "state=URL_ELIGIBLE")
	req.Header.Set("X-Graphql-Client-Library", "pando")
	req.Header.Set("X-Ig-App-Id", "567067343352427")
	req.Header.Set("X-Fb-Request-Analytics-Tags", "{\"network_tags\":{\"product\":\"567067343352427\",\"purpose\":\"none\",\"request_category\":\"graphql\",\"retry_attempt\":\"0\"},\"application_tags\":\"pando\"}")
	req.Header.Set("X-Root-Field-Name", "xfb_pair_generate_text2stickers")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))
	req.Header.Set("X-Ig-Capabilities", "3brTv10=")
	req.Header.Set("X-Fb-Friendly-Name", "IGGenerateText2Stickers")
	req.Header.Set("X-Fb-Http-Engine", "Liger")
	req.Header.Set("X-Fb-Client-Ip", "True")
	req.Header.Set("X-Fb-Server-Cluster", "True")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	respString := string(respBytes)

	var response Response
	err = json.Unmarshal([]byte(strings.ReplaceAll(respString, "1$xfb_pair_generate_text2stickers(bypass_cache:$bypass_cache,caller:$caller,media_type:$media_type,prompt:$prompt,scaling_factor:$scaling_factor)", "xfb_pair_generate_text2stickers")), &response)
	if err != nil {
		panic(err)
	}

	var imgArr []StickerData

	// check if there was an error in generation (can occur when prompted with prohibited terms)
	if response.Data.Data.Error != nil {
		fmt.Println("\n[!] A server-side error occured during generation! Your prompt may contain prohibited terms. Exiting...")
		return
	}

	for i, sticker := range response.Data.Data.Stickers {
		// get the image from the server
		fmt.Printf("\n(*) image %d:\n", i+1)
		imgResponse, err := client.Get(sticker.URL)
		if err != nil {
			panic(err)
		}

		// read the data to add the byte array to the image data slice
		imgBytes, err := io.ReadAll(imgResponse.Body)
		if err != nil {
			panic(err)
		}

		// total jank ! yikes !
		imgArr = append(imgArr,
			StickerData{
				imgBytes,
				"123456" + strings.Split(strings.Split(sticker.URL, "123456")[1], "?")[0],
			})

		// decode the png response
		img, err := png.Decode(bytes.NewReader(imgBytes))

		// print sixel to terminal
		err = sixel.NewEncoder(os.Stdout).Encode(img)
		if err != nil {
			panic(nil)
		}
	}

	// ask user if they want to save
	for {
		fmt.Print("\nsave? (y/n): ")
		saveInput := promptForInput()
		if saveInput == "y" {
			for _, img := range imgArr {
				err := os.WriteFile(img.Filename, img.ImageData, 0666)
				if err != nil {
					panic(err)
				}
			}

			fmt.Println("saved all images to the current directory. see ya later! ✌️")
			return
		} else if saveInput == "n" {
			fmt.Println("sounds good! see ya later! ✌️")
			return
		} else {
			fmt.Println("please only respond with 'y' or 'n'. i'm not the best with the whole alphabet!")
		}
	}

}

func promptForInput() string {
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		panic(err)
	}

	return input
}
