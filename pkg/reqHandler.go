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
	"github.com/ua-parser/uap-go/uaparser"
)

var uaParser *uaparser.Parser

func init() {
	uaParser = uaparser.NewFromSaved() // Loads the default regexes
}

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

func GetFlatHeaders(r *http.Request) map[string]string {
	flatHeaders := make(map[string]string)
	for key, values := range r.Header {
		flatHeaders["header_"+strings.ToLower(key)] = strings.Join(values, ", ")
	}
	return flatHeaders
}

func GetFlatCookies(r *http.Request) map[string]string {
	flatCookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		flatCookies["cookie_"+strings.ToLower(cookie.Name)] = cookie.Value
	}
	return flatCookies
}

func GetHostname(r *http.Request) string {
	host := r.Host
	// Remove the port if it's present
	if colonIndex := strings.IndexByte(host, ':'); colonIndex != -1 {
		host = host[:colonIndex]
	}
	return host
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	ip := r.RemoteAddr
	if colonIndex := strings.LastIndex(ip, ":"); colonIndex != -1 {
		ip = ip[:colonIndex] // Remove port if present
	}
	return ip
}

func GetPort(r *http.Request) string {
	// Check if a port is specified in r.Host
	hostParts := strings.Split(r.Host, ":")
	if len(hostParts) > 1 {
		return hostParts[1] // Port specified in Host header
	}
	// Default port based on scheme
	if r.TLS != nil {
		return "443" // HTTPS
	}
	return "80" // HTTP
}

func GetProtocol(r *http.Request) string {
	// Check if protocol is http or https
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

// Passing `trapConfig` parameter so each instance can handle its own traps independently.
func allHandler(trapConfig []Trap) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trapFlag := "false"
		ua := uaParser.Parse(r.Header.Get("User-Agent"))
		for _, trap := range trapConfig {
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
					if (CheckHeaders(convertMap(behaviour.Request.Headers), r.Header)) && (CheckParams(convertMap(behaviour.Request.Params), params)) {
						trapFlag = "true"
						details := map[string]string{
							"timestamp":          time.Now().Format(time.RFC3339),
							"src_ip":             GetIP(r),
							"dest_port":          GetPort(r),
							"request_method":     r.Method,
							"protocol":           GetProtocol(r),
							"hostname":           GetHostname(r),
							"request_uri":        r.RequestURI,
							"user-agent_browser": ua.UserAgent.Family,
							"user-agent_os":      ua.Os.Family,
							"trapped":            "true",
							"trapped_for":        trap.Basicinfo.Name,
							"user-agent":         r.Header.Get("User-Agent"),
						}
						for key, value := range GetFlatHeaders(r) {
							details[key] = value
						}
						for key, value := range GetFlatCookies(r) {
							details[key] = value
						}
						if ua.Device.Brand != "" {
							details["user-agent_device_brand"] = ua.Device.Brand
						}
						if ua.Device.Model != "" {
							details["user-agent_device_model"] = ua.Device.Model
						}
						if ua.UserAgent.Major != "" || ua.UserAgent.Minor != "" {
							details["user-agent_browser_version"] = fmt.Sprintf("%s.%s", ua.UserAgent.Major, ua.UserAgent.Minor)
						}
						if ua.Os.Major != "" || ua.Os.Minor != "" {
							details["user-agent_os_version"] = fmt.Sprintf("%s.%s", ua.Os.Major, ua.Os.Minor)
						}
						if trap.Basicinfo.RiskRating != "" {
							details["trapped_risk_rating"] = trap.Basicinfo.RiskRating
						}
						if trap.Basicinfo.References != "" {
							details["trapped_references"] = trap.Basicinfo.References
						}
						LogEntry(details)
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
			details := map[string]string{
				"timestamp":          time.Now().Format(time.RFC3339),
				"src_ip":             GetIP(r),
				"dest_port":          GetPort(r),
				"request_method":     r.Method,
				"protocol":           GetProtocol(r),
				"hostname":           GetHostname(r),
				"request_uri":        r.RequestURI,
				"user-agent_browser": ua.UserAgent.Family,
				"user-agent_os":      ua.Os.Family,
				"trapped":            "false",
				"user-agent":         r.Header.Get("User-Agent"),
			}
			for key, value := range GetFlatHeaders(r) {
				details[key] = value
			}
			for key, value := range GetFlatCookies(r) {
				details[key] = value
			}
			if ua.Device.Brand != "" {
				details["user-agent_device_brand"] = ua.Device.Brand
			}
			if ua.Device.Model != "" {
				details["user-agent_device_model"] = ua.Device.Model
			}
			if ua.UserAgent.Major != "" || ua.UserAgent.Minor != "" {
				details["user-agent_browser_version"] = fmt.Sprintf("%s.%s", ua.UserAgent.Major, ua.UserAgent.Minor)
			}
			if ua.Os.Major != "" || ua.Os.Minor != "" {
				details["user-agent_os_version"] = fmt.Sprintf("%s.%s", ua.Os.Major, ua.Os.Minor)
			}
			LogEntry(details)
		}
		trapFlag = "false"
	})
}

func StartHandler(port string, trapConfig []Trap, cert string, key string) {
	r := mux.NewRouter()
	fmt.Println("[~>] Loaded " + strconv.Itoa(len(trapConfig)) + " traps on Port:" + port + ". Let's get the ball rolling!")

	// Pass each port's `trapConfig` directly to `allHandler` to preserve traps on different ports
	r.PathPrefix("/").Handler(allHandler(trapConfig))

	if port == "443" {
		log.Fatal(http.ListenAndServeTLS(":"+port, cert, key, r))
	} else {
		log.Fatal(http.ListenAndServe(":"+port, r))
	}
}
