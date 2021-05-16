package scheduler

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"openvpn-processor/pkg/config"
	"strings"
	"time"
)

func createStructsFromCsv(csvContent [][]string) []vpnServer {
	var vpnServers []vpnServer
	for _, entry := range csvContent {
		server := vpnServer{
			uuid:           uuid.New().String(),
			hostname:       entry[0],
			score:          config.ConvertStringToInt(entry[2]),
			ping:           config.ConvertStringToInt(entry[3]),
			speed:          config.ConvertStringToInt(entry[4]),
			countryLong:    entry[5],
			countryShort:   entry[6],
			numVpnSessions: config.ConvertStringToInt(entry[7]),
			uptime:         config.ConvertStringToInt(entry[8]),
			totalUsers:     config.ConvertStringToInt(entry[9]),
			totalTraffic:   config.ConvertStringToInt(entry[10]),
			enabled:        true,
			createdAt:      time.Now(),
		}

		decodedByteSlice, err := base64.StdEncoding.DecodeString(entry[14])
		if err != nil {
			logger.Warn("an error occurred while decoding conf data, skipping", zap.String("data", entry[0]))
			continue
		}

		decodedConfData := string(decodedByteSlice)
		server.confData = decodedConfData
		for _, line := range strings.Split(decodedConfData, "\n") {
			fields := strings.Fields(line)
			if strings.HasPrefix(line, "remote") {
				server.ip = fields[1]
				server.port = config.ConvertStringToInt(fields[2])
			}

			if strings.HasPrefix(line, "proto") {
				server.proto = fields[1]
			}
		}
		vpnServers = append(vpnServers, server)
	}

	logger.Info("successfully created structs from csv", zap.Int("structsCreated", len(vpnServers)))
	return vpnServers
}

func getCsvContent(vpnGateUrl string) [][]string {
	logger.Info("getting server list from vpngate", zap.String("vpnGateUrl", vpnGateUrl))
	var csvContent [][]string
	resp, err := http.Get(vpnGateUrl)
	if err != nil {
		logger.Error("an error occurred while making GET request", zap.String("vpnGateUrl", vpnGateUrl),
			zap.String("error", err.Error()))
		return nil
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	encodedBody, err := ioutil.ReadAll(resp.Body)
	decodedBody := string(encodedBody)
	if err != nil {
		logger.Error("an error occurred while reading response body", zap.String("vpnGateUrl", vpnGateUrl),
			zap.String("error", err.Error()))
		return nil
	}
	reader := csv.NewReader(strings.NewReader(decodedBody))
	for {
		server, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !strings.HasPrefix(server[0], "*") && !strings.HasPrefix(server[0], "#") {
			csvContent = append(csvContent, server)
		}
	}
	return csvContent
}

func isServerInsertable(ip, proto, confData string, port int, timeoutSeconds int) bool {
	isReachable := true
	timeout := time.Duration(int32(timeoutSeconds)) * time.Second
	_, err := net.DialTimeout(proto, fmt.Sprintf("%s:%d", ip, port), timeout)
	if err != nil {
		isReachable = false
	}

	isUnauthenticated := strings.Contains(confData, "#auth-user-pass")
	return isReachable && isUnauthenticated
}
