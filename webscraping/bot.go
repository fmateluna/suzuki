package webscraping

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// Sirve para obtener la signatura
const KEY_SBSEPC5S = "mHpUcqwUDsLeAR4RPyy-EVDcQnVnWYGVXYuJeYNAU9s"

const KEY_SBSEPC5ACS = "cZPY7mwEsTva7iYodz4ZJhrpR4Dd3WiZHeaBz87FQX0"

type TokenLoginResponse struct {
	Status                   int    `json:"status"`
	SessionJwtToken          string `json:"sessionJwtToken"`
	ProductCodesLicenseCodes []struct {
		ProductCode string `json:"productCode"`
		LicenseCode string `json:"licenseCode"`
	} `json:"productCodesLicenseCodes"`
	LocalSessionPingInterval int `json:"localSessionPingInterval"`
}

type SessionJwtToken struct {
	Cs   string    `json:"CS"`
	Sig  string    `json:"SIG"`
	St   string    `json:"ST"`
	UID  string    `json:"UID"`
	Rd   string    `json:"RD"`
	Cspk string    `json:"CSPK"`
	Rl   []string  `json:"RL"`
	Pk   string    `json:"PK"`
	Vs   string    `json:"VS"`
	Did  string    `json:"DID"`
	Sid  string    `json:"SID"`
	Ts   time.Time `json:"TS"`
}

type BotSuzuki struct {
	User             string
	Pass64           string
	Pass             string
	Cookies          map[string]http.Cookie
	AccountBot       Account
	UserBot          User
	VIN              string
	Sbsepc5acs       string
	Sbsepc5s         string
	TheCookies       string
	SessionJwtToken  SessionJwtToken
	cookieAWSALB     string
	cookieAWSALBCORS string
	cookieJSESSIONID string
	VinObject        VIN
}

func (b *BotSuzuki) Init(vin string) {

	//

	b.Cookies = make(map[string]http.Cookie)
	b.Pass64 = b64.StdEncoding.EncodeToString([]byte(b.Pass))
	if b.login() {
		//log.Println("Easy like sunday morning!")
		b.makeAccountGreatAgain()
		if err := os.Mkdir(vin, os.ModePerm); err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(vin + "/registro.log")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Configura el log para escribir en el archivo
		log.SetOutput(file)

		b.findByVIN(vin)
		b.getCategories()
		b.createJSON()

	} else {
		log.Println("WTF??")
	}
}

func (b *BotSuzuki) createJSON() {

	f, err := os.Create(b.VIN + "\\components.json")

	contentJson, err := json.MarshalIndent(b.VinObject, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(string(contentJson))

	if err2 != nil {
		log.Fatal(err2)
	}

}

func (b *BotSuzuki) GenerateSbsepc5cs() string {
	nuevoSBSEPC5S := SBSEPC5CS{}
	theTime := time.Now().Add(time.Hour * 24 * 2)
	Ts := strings.Split(theTime.Format(time.RFC3339Nano), ".")[0] + ".433Z"

	nuevoSBSEPC5S.Sid = b.SessionJwtToken.Sid
	nuevoSBSEPC5S.Ts = Ts
	nuevoSBSEPC5S.Pk = "SBSEPC5"
	nuevoSBSEPC5S.Rd = b.RandomString(9)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, nuevoSBSEPC5S)
	signedString, _ := token.SignedString([]byte(b.SessionJwtToken.Sig))

	//log.Println(signedString)

	return signedString

}

func (b *BotSuzuki) GenerateSbsepc5s() string {
	nuevoSBSEPC5S := SBSEPC5CS{}
	theTime := time.Now().Add(time.Hour * 24 * 2)
	Ts := strings.Split(theTime.Format(time.RFC3339Nano), ".")[0] + ".433Z"

	nuevoSBSEPC5S.Sid = b.SessionJwtToken.Sid
	nuevoSBSEPC5S.Ts = Ts
	nuevoSBSEPC5S.Pk = "SBSEPC5"
	nuevoSBSEPC5S.Rd = b.RandomString(9)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, nuevoSBSEPC5S)
	signedString, _ := token.SignedString([]byte(KEY_SBSEPC5S))

	//log.Println(signedString)

	return signedString
}

func (b *BotSuzuki) makeAccountGreatAgain() {

	req, err := http.NewRequest("POST", "https://suzuki.snaponepc.com/epc-services/auth/account", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "no-cache,no-store")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "JSESSIONID="+b.cookieJSESSIONID+"; AWSALB="+b.cookieAWSALB+"; AWSALBCORS="+b.cookieAWSALBCORS)
	req.Header.Set("Sbsepc5cs", b.GenerateSbsepc5cs())
	req.Header.Set("Sbsepc5s", b.GenerateSbsepc5s())
	req.Header.Set("Expires", "0")
	req.Header.Set("Origin", "https://suzuki.snaponepc.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	b.takeCookies(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	account := Account{}
	accountResponse := string(body)
	//log.Println(accountResponse)
	decoder := json.NewDecoder(strings.NewReader(accountResponse))
	decoder.Decode(&account)
	b.AccountBot = account
	defer resp.Body.Close()

}

func (b *BotSuzuki) CreateSBSEPC5ACS() string {

	nuevoSBSEPC5S := SBSEPC5CS{}
	theTime := time.Now().Add(time.Hour * 24 * 2)
	Ts := strings.Split(theTime.Format(time.RFC3339Nano), ".")[0] + ".433Z"
	nuevoSBSEPC5S.Ts = Ts
	nuevoSBSEPC5S.Pk = "SBSEPC5"
	nuevoSBSEPC5S.Rd = b.RandomString(9)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, nuevoSBSEPC5S)
	signedString, _ := token.SignedString([]byte(KEY_SBSEPC5ACS))

	//log.Println(signedString)

	return signedString
}

func (b *BotSuzuki) login() bool {
	params := url.Values{}
	params.Add("user", b.User)
	params.Add("password", b.Pass64)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://suzuki.snaponepc.com/epc-services/auth/login", body)
	if err != nil {
		// handle err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "no-cache,no-store")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "AWSALB=MPT4nx8NCqNIjcpuRKymIHf3r50ptt26fie7Vbiv0lH+Chh4cKH3ENBQIRAnZ/njw/AH23GZo0J3B+tivXZNHcFA0MYPv67RdWNcIaQOYKisRdrT16PFVBnsONXs; AWSALBCORS=MPT4nx8NCqNIjcpuRKymIHf3r50ptt26fie7Vbiv0lH+Chh4cKH3ENBQIRAnZ/njw/AH23GZo0J3B+tivXZNHcFA0MYPv67RdWNcIaQOYKisRdrT16PFVBnsONXs")
	req.Header.Set("Expires", "0")
	req.Header.Set("Origin", "https://suzuki.snaponepc.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Sbsepc5acs", b.CreateSBSEPC5ACS())
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()

	if err != nil {
		return false
	} else {
		//log.Println("PONG!! ", resp.StatusCode)
		if resp.StatusCode == 401 || resp.StatusCode == http.StatusForbidden {
			return false
		} else {
			bodyBytes, err := io.ReadAll(resp.Body)
			tokenResponse := TokenLoginResponse{}
			jsonContent := string(bodyBytes)
			decoder := json.NewDecoder(strings.NewReader(jsonContent))
			decoder.Decode(&tokenResponse)
			b.takeCookies(resp)
			if err != nil {
				return false
			}
			base64 := string(tokenResponse.SessionJwtToken)

			payload64 := strings.Split(base64, ".")[1] + "="

			payloadJson, err := b64.StdEncoding.DecodeString(payload64)
			//Println(string(payloadJson))

			sessionJwtToken := SessionJwtToken{}
			decoder = json.NewDecoder(strings.NewReader(string(payloadJson)))
			decoder.Decode(&sessionJwtToken)

			b.SessionJwtToken = sessionJwtToken

			return true
		}
	}

}

func (b *BotSuzuki) takeCookies(resp *http.Response) {

	if resp.Request != nil {
		//log.Println("Seteando REQUEST Cookies()  ", len(resp.Request.Cookies()))
		for _, cookie := range resp.Request.Cookies() {
			//log.Println("[+] Cookies ", cookie.Name, " =", cookie.Value)
			b.Cookies[cookie.Name] = *cookie
		}
	}

	//log.Println("LEYENDO Response Cookies() ", len(resp.Cookies()))
	for _, cookie := range resp.Cookies() {
		//log.Println("[!] Cookies ", cookie.Name, " =", cookie.Value)

		if cookie.Name == "AWSALB" {
			b.cookieAWSALB = cookie.Value
		}
		if cookie.Name == "AWSALBCORS" {
			b.cookieAWSALBCORS = cookie.Value
		}
		if cookie.Name == "JSESSIONID" {
			b.cookieJSESSIONID = cookie.Value
		}
		b.Cookies[cookie.Name] = *cookie
	}
}

func (b *BotSuzuki) downloadImagen(imageId string) {

	url := "https://suzuki.snaponepc.com/epc-services/datasets/" + b.VinObject.DatasetID + "/pages/images/" + imageId

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Amg", "f1360714-40d9-4416-b744-e3ec8d98ac2f")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Cookie", "JSESSIONID="+b.cookieJSESSIONID+"; AWSALB="+b.cookieAWSALB+"; AWSALBCORS="+b.cookieAWSALBCORS)
	req.Header.Set("Sbsepc5cs", b.GenerateSbsepc5cs())
	req.Header.Set("Sbsepc5s", b.GenerateSbsepc5s())
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	file, err := os.Create(b.VIN + "\\" + imageId + ".png")

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	defer resp.Body.Close()
}

func (b *BotSuzuki) findByVIN(vin string) {
	b.VIN = vin
	base64 := b64.URLEncoding.EncodeToString([]byte(b.returnFirma()))
	url := "https://suzuki.snaponepc.com/epc-services/equipment/search?q=" + b.VIN + "&fr=" + base64 + "&es=true"

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Amg", "f1360714-40d9-4416-b744-e3ec8d98ac2f")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Cookie", "JSESSIONID="+b.cookieJSESSIONID+"; AWSALB="+b.cookieAWSALB+"; AWSALBCORS="+b.cookieAWSALBCORS)
	req.Header.Set("Sbsepc5cs", b.GenerateSbsepc5cs())
	req.Header.Set("Sbsepc5s", b.GenerateSbsepc5s())
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	vinResult := SearchVIN{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&vinResult)

	b.VinObject = vinResult.VinSearchResults[0].Vins[0]

	defer resp.Body.Close()
}

func (b *BotSuzuki) getCategories() {

	url := "https://suzuki.snaponepc.com/epc-services/datasets/" + b.VinObject.DatasetID + "/navigations/" + b.VinObject.SerializedPath + "/filterRequest/" + b.generateFilterRequest()

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Amg", "f1360714-40d9-4416-b744-e3ec8d98ac2f")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Cookie", "JSESSIONID="+b.cookieJSESSIONID+"; AWSALB="+b.cookieAWSALB+"; AWSALBCORS="+b.cookieAWSALBCORS)
	req.Header.Set("Sbsepc5cs", b.GenerateSbsepc5cs())
	req.Header.Set("Sbsepc5s", b.GenerateSbsepc5s())
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	responseFilter := ResponseFilters{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&responseFilter)

	for _, categoria := range responseFilter.Children.ChildNodes {

		log.Println(b.VIN + "\\ Category : " + categoria.Name)

		myCategoria := Category{}
		myCategoria.Filtered = categoria.Filtered
		myCategoria.HasNotes = categoria.HasNotes
		myCategoria.ID = categoria.ID
		myCategoria.ImageID = categoria.ImageID
		myCategoria.LeafNode = categoria.LeafNode
		myCategoria.Name = categoria.Name
		myCategoria.SerializedPath = categoria.SerializedPath
		b.downloadImagen(categoria.ImageID)
		b.subCategory(&myCategoria)
		b.VinObject.Categories = append(b.VinObject.Categories, myCategoria)
	}

	defer resp.Body.Close()
}

func (b *BotSuzuki) subCategory(subCategory *Category) {

	url := "https://suzuki.snaponepc.com/epc-services/datasets/" + b.VinObject.DatasetID + "/navigations/" + subCategory.SerializedPath + "/filterRequest/" + b.generateFilterRequest()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Amg", "f1360714-40d9-4416-b744-e3ec8d98ac2f")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Cookie", "JSESSIONID="+b.cookieJSESSIONID+"; AWSALB="+b.cookieAWSALB+"; AWSALBCORS="+b.cookieAWSALBCORS)
	req.Header.Set("Sbsepc5cs", b.GenerateSbsepc5cs())
	req.Header.Set("Sbsepc5s", b.GenerateSbsepc5s())
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"110\", \"Not A(Brand\";v=\"24\", \"Google Chrome\";v=\"110\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	responseFilter := ResponseFilters{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&responseFilter)

	for _, pieza := range responseFilter.Children.ChildNodes {
		pieza.Parts = b.getPieza(subCategory, pieza, pieza.ImageID)
		b.downloadImagen(pieza.ImageID)
		log.Println(b.VIN + "\\ - SubCategory : " + pieza.Name)
		subCategory.SubCategory = append(subCategory.SubCategory, pieza)
	}

	defer resp.Body.Close()

}

func (b *BotSuzuki) getPieza(padre *Category, pieza SubCategory, imageID string) ResponsePiezas {
	url := "https://suzuki.snaponepc.com/epc-services/datasets/" + b.VinObject.DatasetID + "/pages/parts/" + pieza.SerializedPath + "/filterRequest/" + b.generateFilterRequest() + "?imageId=" + imageID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Amg", "f1360714-40d9-4416-b744-e3ec8d98ac2f")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "es-ES,es;q=0.9")
	req.Header.Set("Cache-Control", "no-cache,no-store")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Expires", "0")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://suzuki.snaponepc.com/epc/")
	req.Header.Set("Cookie", "JSESSIONID="+b.cookieJSESSIONID+"; AWSALB="+b.cookieAWSALB+"; AWSALBCORS="+b.cookieAWSALBCORS)
	req.Header.Set("Sbsepc5cs", b.GenerateSbsepc5cs())
	req.Header.Set("Sbsepc5s", b.GenerateSbsepc5s())
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"111\", \"Not(A:Brand\";v=\"8\", \"Chromium\";v=\"111\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Sec-Ch-Ua-Platform-Version", "\"14.0.0\"")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	responsePieza := ResponsePiezas{}
	jsonContent := string(body)
	decoder := json.NewDecoder(strings.NewReader(jsonContent))
	decoder.Decode(&responsePieza)

	//fmt.Println("MODELO;CATEGORIA;SUBCATEGORIA;PIEZA;ID PARTE;PARTE;DESCRIPCION")

	for _, pageImage := range responsePieza.PageImages {
		b.downloadImagen(pageImage.ImageID)
		//subPart := b.getPieza(padre, pieza, pageImage.ImageID)
		//responsePieza.SubParts = append(responsePieza.SubParts, subPart)
	}

	return responsePieza
}

func (b *BotSuzuki) generateFilterRequest() string {
	filterRequest := fmt.Sprintf("jobId=1|dataSetId=%s|manualFiltersEnabled=false|equipmentRefId=%s|currentVin=%s|currentVinBusRegRef=%s|vinBusRegRef=%s|einId=%s|filtersEnabled=true|locale=en-US|busReg=SUZ|LA|userId=%s",
		b.VinObject.DatasetID,
		b.VinObject.EquipmentRefID,
		b.VinObject.FormattedVin,
		b.VinObject.BusinessRegion,
		b.VinObject.BusinessRegion,
		b.VinObject.EinID,
		b.AccountBot.UserDetails.UserID,
	)
	return b64.StdEncoding.EncodeToString([]byte(filterRequest))
}
