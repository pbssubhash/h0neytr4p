package h0neytr4p

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ryanuber/go-glob"
)

var trapConfigGlobal []Trap
var trapFlag = "false"

func convertMap(mapInterface map[string]interface{}) map[string]string {
	mapString := make(map[string]string)
	for key, value := range mapInterface {
		strKey := fmt.Sprintf("%v", key)
		strValue := fmt.Sprintf("%v", value)
		mapString[strKey] = strValue
	}
	return mapString
}

func match(first string, second string) bool {
	return glob.Glob(first, second)
}

func CheckHeaders(ruleHeaders map[string]string, requestHeaders http.Header) bool {
	for k, v := range ruleHeaders {
		if !match(v, requestHeaders.Get(k)) {
			return false
		}
	}
	return true
}
func CheckParams(ruleParams map[string]string, requestParams map[string]string) bool {
	for k, v := range ruleParams {
		if !match(v, requestParams[k]) {
			return false
		}
	}
	return true
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func allHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, trap := range trapConfigGlobal {
			for _, behaviour := range trap.Behaviour {
				params := make(map[string]string)
				if (behaviour.Request.Method == r.Method) && (match(behaviour.Request.URL, r.URL.Path)) {
					if r.Method == "GET" {
						if len(strings.Split(r.RequestURI, "?")) > 1 {
							for _, k := range strings.Split(strings.Split(r.RequestURI, "?")[1], "&") {
								split := strings.Split(k, "=")
								if len(split) == 2 {
									params[split[0]] = split[1]
								} else {
									params[split[0]] = ""
								}
							}
						}
					} else if r.Method == "POST" || r.Method == "PUT" || r.Method == "DELETE" {
						if len(strings.Split(r.RequestURI, "?")) > 1 {
							for _, k := range strings.Split(strings.Split(r.RequestURI, "?")[1], "&") {
								split := strings.Split(k, "=")
								if len(split) == 2 {
									params[split[0]] = split[1]
								} else {
									params[split[0]] = ""
								}
							}
						}
						if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
							if err := r.ParseMultipartForm(0); err != nil {
								fmt.Println("Error with Params")
								fmt.Println(err.Error())
							}
							for key, values := range r.PostForm {
								params[key] = strings.Join(values, "|")
							}
						} else if r.Header.Get("Content-Type") == "text/plain" {
							bodyBytes, err := ioutil.ReadAll(r.Body)
							if err != nil {
								log.Fatal(err)
							}
							bodyString := string(bodyBytes)
							for _, k := range strings.Split(bodyString, "&") {
								split := strings.Split(k, "=")
								if len(split) == 2 {
									params[split[0]] = split[1]
								} else {
									params[split[0]] = ""
								}
							}
						} else if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
							if err := r.ParseForm(); err != nil {
								fmt.Println("Error with Params")
							}
							for key, values := range r.PostForm {
								params[key] = strings.Join(values, "|")
							}
						}
					}
					fmt.Println(params)
					trapFlag = "false"
					if (CheckHeaders(convertMap(behaviour.Request.Headers), r.Header)) && (CheckParams(convertMap(behaviour.Request.Params), params)) {
						trapFlag = "true"
						LogEntry(LogDetails{SourceIP: GetIP(r),
							UserAgent: r.Header.Get("User-Agent"),
							Timestamp: time.Now().Format(time.RFC3339),
							Path:      r.RequestURI, Trapped: "true",
							TrappedFor: trap.Basicinfo.Name,
							RiskRating: trap.Basicinfo.RiskRating,
							References: trap.Basicinfo.References})
						// Writing Response according to trap
						responseHeaders := convertMap(behaviour.Response.Headers)
						for key, value := range responseHeaders {
							w.Header().Set(key, value)
						}
						w.WriteHeader(behaviour.Response.Statuscode)
						if behaviour.Response.Type == "file" {
							content, err := ioutil.ReadFile(behaviour.Response.Body)
							if err != nil {
								fmt.Println("[RESPONSE-ERROR]: Unable to read file")
							} else {
								_, err = w.Write(content)
								if err != nil {
									fmt.Println("Unable to write content")
								}
							}
						} else {
							w.Write([]byte(behaviour.Response.Body))
						}
						// End Writing Response
					}
				}
			}
		}
		if trapFlag != "true" {
			LogEntry(LogDetails{SourceIP: GetIP(r),
				UserAgent: r.Header.Get("User-Agent"),
				Timestamp: time.Now().Format(time.RFC3339),
				Path:      r.RequestURI, Trapped: "false",
				TrappedFor: "",
				RiskRating: "",
				References: ""})
		}
		trapFlag = "false"
	})
}

func StartHandler(port string, trapConfig []Trap, cert string, key string) {
	r := mux.NewRouter()
	fmt.Println("[~>] Loaded " + strconv.Itoa(len(trapConfig)) + " traps on Port:" + port + ". Let's get the ball rolling!")
	trapConfigGlobal = trapConfig
	// r.HandleFunc("/", allHandler)
	r.PathPrefix("/").Handler(allHandler())
	if port == "443" {
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, r))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}
