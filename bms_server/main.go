package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/ronnie-w/bms-server/inventories"
	"github.com/ronnie-w/bms-server/notifications"
	"github.com/ronnie-w/bms-server/router"
	"github.com/ronnie-w/bms-server/sales"
	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/staff"
	"github.com/ronnie-w/bms-server/utils"
	webview "github.com/webview/webview_go"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	/*var (
		f, err = os.OpenFile("crash.dump", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

		logger = func(strerr string, err error) {
			f.Write([]byte(time.Now().String() + strerr + err.Error()))
			os.Exit(3)
		}
	)
	defer f.Close()

	ifas, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}

	var mac []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			mac = append(mac, a)
		}
	}

	key := []byte("Hs&x+FS_LmT@Y^F#Hc-5hbJ'k_{$Lwp9")
	cipherText, err := os.ReadFile("licence")
	if err != nil || len(cipherText) == 0 {
		logger("Fatal Error! Quiting (xe1)...", err)
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		logger("Fatal Error! Quiting (xe2)...", err)
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		logger("Fatal Error! Quiting (xe3)...", err)
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		logger("Fatal Error! Quiting (xe4)...", err)
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	gcmHashArr, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		logger("!xoxo!", err)
	}

	var (
		collectHash [][64]byte
		hashArr     []byte
	)
	for _, m := range mac {
		collectHash = append(collectHash, sha512.Sum512([]byte(m)))
	}

	for _, ch := range collectHash {
		hashArr = append(hashArr, ch[:]...)
	}

	if !bytes.Equal(gcmHashArr, hashArr) {
		logger("!xoxo!", nil)
	} else {
		fmt.Println("xoxo")
	}*/

	if err := schema.CreateTableQuery("tabs", []string{"tab_name varchar(255) not null"}); err == nil {
		defaultTabs := []string{"sales", "customers"}
		for i := 0; i < len(defaultTabs); i++ {
			_ = schema.InsertQuery("tabs", []interface{}{defaultTabs[i]})
		}

		inventories.CreateInventoryTable()
		sales.CreateSalesTable()
		notifications.CreateNotificationsTable()
		staff.CreateStaffTable()
		sales.CreateCustomerTable()
		utils.CreateSettingsTable()
	}
}

func main() {
	guiClosed := make(chan bool)
	go func() {
		w := webview.New(false)
		w.SetTitle("Beazy")
		w.SetSize(960, 540, webview.HintNone)

		defer w.Destroy()
		w.Navigate("http://localhost:8000")
		w.Run()

		guiClosed <- true
	}()

	go func() {
		var validator handlers.OriginValidator = func(s string) bool {
			return strings.Contains(s, "localhost") || strings.Contains(s, "127.0.0.1")
		}

		methods := handlers.AllowedMethods([]string{"POST", "GET"})
		origin := handlers.AllowedOriginValidator(validator)
		credentials := handlers.AllowCredentials()

		logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println(err)
		}

		server := &http.Server{
			Addr:    ":8000",
			Handler: handlers.LoggingHandler(logFile, handlers.CORS(methods, origin, credentials)(handlers.CompressHandlerLevel(router.Routes(), 9))),
		}

		if err := server.ListenAndServe(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case <-guiClosed:
			os.Exit(0)
		}
	}
}
