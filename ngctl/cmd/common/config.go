package common

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type IPAndPort struct {
	// define both members as strings, since we only need them to fill config files
	IP          string
	Port        string
	AgentPort   string
	ServiceType string
}

// Genrate config files for metad services only
func GenerateMetadConfigFile(defaultConfigFile string, resultConfigFilePath string, destination IPAndPort, metadAddrs []IPAndPort) (err error) {
	file, err := os.Open(defaultConfigFile)
	if err != nil {
		return fmt.Errorf("failed to open default config file: %v", err)
	}
	defer file.Close()

	// Create a new file to write the modified config
	newFile, err := os.Create(resultConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to create new config file: %v", err)
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// The following configs need to be updated from the default: --meta_server_addrs, --local_ip, --port
		// Other paramters are not updated here
		// Update the necessary configurations in the line
		if strings.Contains(line, "--meta_server_addrs") {
			line = "--meta_server_addrs="
			metaAddrsStr := ""
			for _, addr := range metadAddrs {
				metaAddrsStr += addr.IP + ":" + fmt.Sprint(addr.Port) + ","
			}
			line = line + strings.TrimSuffix(metaAddrsStr, ",")
		} else if strings.Contains(line, "--port") {
			line = "--port=" + fmt.Sprint(destination.Port)
		} else if strings.Contains(line, "--local_ip") {
			line = "--local_ip=" + destination.IP
		}
		// Write the modified line to the new config file
		_, err := newFile.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line to new config file: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading default config file: %v", err)
	}

	return nil
}

// Generate config files for a single graphd or storage service
func GenerateConfigFile(serviceType string, defaultConfigFile string, resultConfigFilePath string, destination IPAndPort, metadAddrs []IPAndPort) (err error) {
	// The defulat config files are stored in the ngctl/etc/ folder:
	// nebula-graphd.conf.default  nebula-metad.conf.default  nebula-storaged.conf.default
	// This function updates the IPs and ports in these config files. Leave all the rest as default.
	if serviceType != "graphd" && serviceType != "storaged" {
		return fmt.Errorf("only generating config files for graphd or storaged services in this function")
	}

	// Read the default config file
	file, err := os.Open(defaultConfigFile)
	if err != nil {
		return fmt.Errorf("failed to open default config file: %v", err)
	}
	defer file.Close()

	// Create a new file to write the modified config
	newFile, err := os.Create(resultConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to create new config file: %v", err)
	}
	defer newFile.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// The following configs need to be updated from the default: --meta_server_addrs, --local_ip, --port
		// Other paramters are not updated here
		// Update the necessary configurations in the line
		if strings.Contains(line, "--meta_server_addrs") {
			line = "--meta_server_addrs="
			metaAddrsStr := ""
			for _, addr := range metadAddrs {
				metaAddrsStr += addr.IP + ":" + fmt.Sprint(addr.Port) + ","
			}
			line = line + strings.TrimSuffix(metaAddrsStr, ",")
		} else if strings.Contains(line, "--port") {
			line = "--port=" + fmt.Sprint(destination.Port)
		} else if strings.Contains(line, "--local_ip") {
			line = "--local_ip=" + destination.IP
		}
		// Write the modified line to the new config file
		_, err := newFile.WriteString(line + "\n")
		if err != nil {
			return fmt.Errorf("failed to write line to new config file: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading default config file: %v", err)
	}

	return nil
}

func BackupFile(srcPath string, dstPath string) (err error) {
	if src, err := os.ReadFile(srcPath); err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	} else if err := os.WriteFile(dstPath, src, 0644); err != nil {
		return fmt.Errorf("failed to write the copied file: %v", err)
	}
	return nil
}
