package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/console/cache"
	"github.com/vesoft-inc/nebula-ng-tools/console/printer"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

var metaClient meta.Client

func metaClientInit() error {
	if metaClient != nil {
		return nil
	}
	cacheSession, err := cache.LoadMetaSession()
	if err != nil {
		return fmt.Errorf("load meta session failed: %s", err)
	}
	metaClient, err = meta.NewMetaClient(cacheSession.Address)
	if err != nil {
		return err
	}
	return nil
}

func metaClientClose() {
	if metaClient != nil {
		metaClient.Close()
	}
}

func metaConsoleError(message string, err string) error {
	if err != "" {
		return fmt.Errorf("%s, err: %s", message, err)
	} else {
		return fmt.Errorf("%s", message)
	}
}

var rootCmd = &cobra.Command{
	Use:   "meta-console",
	Short: "Execute meta command in cli mode.",
	Long: `Execute meta command in cli mode. Use 'meta-console -h' to see usage.
	**Notice:** You should login meta server first`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// ***************** to connect/disconnect command *****************

type LoginFlags struct {
	Metas string // meta server address list separated by comma like "xx:xx,xx:xx"
	Addr  string
	Port  uint32
	User  string
	Pass  string
}

var loginFlags LoginFlags

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login meta server.",
	Long:  `login meta server --addr [ip] --port [port] --user [user] --password [password]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var address string
		if loginFlags.Metas == "" {
			address = fmt.Sprintf("%s:%d", loginFlags.Addr, loginFlags.Port)
		} else {
			metas := strings.Split(loginFlags.Metas, ",")
			if len(metas) == 0 {
				return fmt.Errorf("meta server address is empty")
			}
			address = loginFlags.Metas
		}
		_, err := meta.NewMetaClient(address)
		if err != nil {
			return err
		}
		cache.SaveMetaSession(address)
		// TODO should login?
		// metaclient.Login(loginFlags.User, loginFlags.Pass)

		fmt.Println("[Warning] Login meta is not implemented yet.")
		fmt.Println("You can do operations without login now.")

		return nil
	},
}

// ***************** to manage cluster *****************
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Process cluster command",
	Long:  `Execute cluster command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type CreateClusterFlags struct {
	ClusterName string
	Replica     int
	Zones       []string
	IfNotExists bool
}

var createClusterFlags CreateClusterFlags

var metaCreateClusterCmd = &cobra.Command{
	Use:   "create",
	Short: "Create cluster in meta server.",
	Long:  `meta-console cluster create --cluster [clustername] --replica [replica] --zones [zone1,zone2,...] --if_not_exists`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := createClusterFlags
		if flags.ClusterName == "" {
			return metaConsoleError("Cluster name is empty", "")
		}
		if flags.Replica == 0 {
			return metaConsoleError("Replica number is invalid", "")
		}
		req := meta.NewCreateClusterReq(flags.ClusterName, flags.Replica, flags.Zones)
		resp, err := metaClient.CreateCluster(req)
		if err != nil {
			return metaConsoleError("Create cluster failed", err.Error())
		}
		if resp.GetErrorCode() != nebula.ErrorSuccessfulCompletion {
			return metaConsoleError("Create cluster failed", resp.GetErrorMsg())
		}
		fmt.Println("Create cluster successfully.")
		return nil
	},
}

type InitClusterFlags struct {
	Cluster string
}

var initClusterFlags InitClusterFlags

var metaInitClusterCmd = &cobra.Command{
	Use:   "init",
	Short: "Init cluster storage part.",
	Long:  `meta-console cluster init --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := initClusterFlags.Cluster
		if cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewInitClusterReq(cluster)

		resp, err := metaClient.InitCluster(req)
		if err != nil {
			return metaConsoleError("Init cluster failed", err.Error())
		}
		if resp.Code != nebula.ErrorSuccessfulCompletion {
			return metaConsoleError("Init cluster failed", resp.GetErrorMsg())
		}
		fmt.Println("Init cluster successfully.")
		return nil
	},
}

type ShowClusterFlags struct {
	Cluster string
}

var showClusterFlags ShowClusterFlags

var metaShowClusterCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cluster, show all if no cluster name specified.",
	Long:  `nebula-meta cluster show --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := showClusterFlags.Cluster
		req := meta.NewShowClusterReq(cluster)
		resp, err := metaClient.ShowCluster(req)
		if err != nil {
			return metaConsoleError("Show cluster failed", err.Error())
		}
		if resp.GetErrorCode() != nebula.ErrorSuccessfulCompletion {
			return metaConsoleError("Show cluster failed", resp.GetErrorMsg())
		}
		header := []string{"cluster id", "cluster name", "replica", "zones"}
		data := make([][]string, 0)
		for _, s := range resp.Clusters {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.ClusterId))
			row = append(row, fmt.Sprintf("%s", s.ClusterName))
			row = append(row, fmt.Sprintf("%d", s.Replica))
			row = append(row, fmt.Sprintf("%s", strings.Join(s.Zones, ",")))
			data = append(data, row)
		}
		// printer.FormatTable(headers []string, data [][]string)
		fmt.Println(printer.FormatTable(header, data))
		return nil
	},
}

// ***************** to manage service *****************
var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Process service command",
	Long:  `Execute service command in cli mode.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type AddServiceFlags struct {
	ServiceType string
	IP          string
	Port        uint32
	Cluster     string
}

var addServiceFlags AddServiceFlags

var metaAddServiceCmd = &cobra.Command{
	Use:   "add",
	Short: `Add service into assigned cluster.`,
	Long:  `meta-console service add --type [graph|storage] --ip [ip] --port [port] --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := addServiceFlags
		var serviceType meta.ServiceType
		if flags.ServiceType == "graph" {
			serviceType = meta.ServiceTypeGraphd
		} else if flags.ServiceType == "storage" {
			serviceType = meta.ServiceTypeStoraged
		} else {
			return metaConsoleError("service type is not correct, valid value is graph or storage", "")
		}
		if flags.IP == "" {
			return metaConsoleError("service ip is empty", "")
		}

		if flags.Cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewAddServiceReq(flags.IP, flags.Port, serviceType, flags.Cluster)
		resp, err := metaClient.AddService(req)
		if err != nil {
			return metaConsoleError("Add service failed", err.Error())
		}
		if resp.GetErrorCode() != nebula.ErrorSuccessfulCompletion {
			return metaConsoleError("Add service failed", resp.GetErrorMsg())
		}
		fmt.Println("Add service successfully.")
		return nil
	},
}

type ShowServiceFlags struct {
	Cluster string
}

var showServiceFlags ShowServiceFlags

var metaShowServiceCmd = &cobra.Command{
	Use:   "show",
	Short: "Show service in cluster.",
	Long:  `meta-console service show --cluster [clustername]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return metaClientInit()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		metaClientClose()
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster := showServiceFlags.Cluster
		if cluster == "" {
			return metaConsoleError("cluster name is empty", "")
		}
		req := meta.NewShowServiceReq(cluster)

		resp, err := metaClient.ShowService(req)
		if err != nil {
			return metaConsoleError("Show service failed", err.Error())
		}
		if resp.Code != nebula.ErrorSuccessfulCompletion {
			return metaConsoleError("Show service failed", resp.GetErrorMsg())
		}
		header := []string{"service id", "service type", "host", "port"}
		data := make([][]string, 0)
		for _, s := range resp.Services {
			row := make([]string, 0)
			row = append(row, fmt.Sprintf("%d", s.ServiceId))
			if s.ServiceType == meta.ServiceTypeGraphd {
				row = append(row, "graphd")
			} else {
				row = append(row, "storaged")
			}
			row = append(row, fmt.Sprintf("%s", s.Host))
			row = append(row, fmt.Sprintf("%d", s.Port))
			data = append(data, row)
		}
		fmt.Println(printer.FormatTable(header, data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(metaCreateClusterCmd)
	clusterCmd.AddCommand(metaInitClusterCmd)
	clusterCmd.AddCommand(metaShowClusterCmd)

	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(metaAddServiceCmd)
	serviceCmd.AddCommand(metaShowServiceCmd)

	loginCmd.Flags().StringVarP(&loginFlags.Addr, "addr", "a", "", "meta server address")
	loginCmd.Flags().Uint32VarP(&loginFlags.Port, "port", "", 0, "meta server port")
	loginCmd.Flags().StringVarP(&loginFlags.User, "user", "u", "", "user name")
	loginCmd.Flags().StringVarP(&loginFlags.Pass, "password", "p", "", "password")
	loginCmd.Flags().StringVarP(&loginFlags.Metas, "metas", "m", "", "meta server address list separated by comma like \"xx:xx,xx:xx\"")

	metaCreateClusterCmd.Flags().StringVarP(&createClusterFlags.ClusterName, "cluster", "c", "", "cluster name")
	metaCreateClusterCmd.Flags().IntVarP(&createClusterFlags.Replica, "replica", "r", 0, "replica number")
	metaCreateClusterCmd.Flags().StringArrayVarP(&createClusterFlags.Zones, "zones", "z", []string{}, "zones")
	metaCreateClusterCmd.Flags().BoolVarP(&createClusterFlags.IfNotExists, "if_not_exists", "", false, "if not exists")

	metaInitClusterCmd.Flags().StringVarP(&initClusterFlags.Cluster, "cluster", "c", "", "cluster name")

	metaShowClusterCmd.Flags().StringVarP(&showClusterFlags.Cluster, "cluster", "c", "", "cluster name")

	metaAddServiceCmd.Flags().StringVarP(&addServiceFlags.ServiceType, "type", "t", "", "service type")
	metaAddServiceCmd.Flags().StringVarP(&addServiceFlags.IP, "ip", "i", "", "service ip")
	metaAddServiceCmd.Flags().Uint32VarP(&addServiceFlags.Port, "port", "p", 0, "service port")
	metaAddServiceCmd.Flags().StringVarP(&addServiceFlags.Cluster, "cluster", "c", "", "cluster name")

	metaShowServiceCmd.Flags().StringVarP(&showServiceFlags.Cluster, "cluster", "c", "", "cluster name")
}

func main() {
	rootCmd.Execute()
}
