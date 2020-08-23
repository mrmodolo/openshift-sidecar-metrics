//Golang basics - fetch JSON from an API
//https://blog.alexellis.io/golang-json-api-client/

//How To Trust Extra CA Certs In Your Go App
//https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/

package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//http://json2struct.mervine.net/
type PodMetrics struct {
	APIVersion string `json:"apiVersion"`
	Containers []struct {
		Name  string `json:"name"`
		Usage struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	} `json:"containers"`
	Kind     string `json:"kind"`
	Metadata struct {
		CreationTimestamp string `json:"creationTimestamp"`
		Name              string `json:"name"`
		Namespace         string `json:"namespace"`
		SelfLink          string `json:"selfLink"`
	} `json:"metadata"`
	Timestamp string `json:"timestamp"`
	Window    string `json:"window"`
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

//${APISERVER}/apis/metrics.k8s.io/v1beta1/namespaces/${NAMESPACE}/pods/${POD}/
func build_api_url(apiserver string, namespace string, pod string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s/apis/metrics.k8s.io/v1beta1/namespaces/%s/pods/%s", apiserver, namespace, pod)
	return b.String()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func get_token_from_file(filename string) string {
	dat, err := ioutil.ReadFile(filename)
	check(err)
	return strings.TrimSuffix(string(dat), "\n")
}

func get_http_transport(localCertFile string) *http.Transport {

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := ioutil.ReadFile(localCertFile)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", localCertFile, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only")
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		RootCAs: rootCAs,
	}

	tr := &http.Transport{TLSClientConfig: config}

	return tr
}

func main() {

	localCertFile := getEnv("CERT_FILE", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	namespace := getEnv("NAMESPACE", "")
	pod := getEnv("POD", "")
	apiserver := getEnv("APISERVER", "https://kubernetes.default")
	token_file := getEnv("TOKEN_FILE", "/var/run/secrets/kubernetes.io/serviceaccount/token")
	token := get_token_from_file(token_file)

	url := build_api_url(apiserver, namespace, pod)

	spaceClient := http.Client{
		Timeout:   time.Second * 10, // Timeout after 2 seconds
		Transport: get_http_transport(localCertFile),
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "curl")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	podMetrics := PodMetrics{}
	containers := PodMetrics{}.Containers
	metadata := PodMetrics{}.Metadata

	jsonErr := json.Unmarshal(body, &podMetrics)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	containers = podMetrics.Containers
	metadata = podMetrics.Metadata

	fmt.Printf("APIVersion: %s\n", podMetrics.APIVersion)

	for _, container := range containers {
		name := metadata.Name
		usage := container.Usage
		fmt.Printf("Pod name: %s\n", name)
		fmt.Println("Usage")
		fmt.Printf("  CPU: %s, Memory: %s\n", usage.CPU, usage.Memory)
	}
}
