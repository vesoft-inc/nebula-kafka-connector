package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/console/cache"
	"github.com/vesoft-inc/nebula-ng-tools/console/printer"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/meta"
)

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Execute meta command in cli mode.",
	Long: `Execute meta command in cli mode. Use 'meta-console -h' to see usage.
	**Notice:** You should login meta server first.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type LoginFlags struct {
	Metas string // meta server address list separated by comma like "xx:xx,xx:xx"
	Addr  string
	Port  uint32
	User  string
	Pass  string
}

var loginFlags LoginFlags
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

var metaLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login meta server.",
	Long:  `login meta server --addr [ip] --port [port] --user [user] --password [password]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginFlags.Metas == "" {
			address := fmt.Sprintf("%s:%d", loginFlags.Addr, loginFlags.Port)
			_, err := meta.NewMetaClient(address)
			if err != nil {
				return err
			}
			// TODO should login?
			// metaclient.Login(loginFlags.User, loginFlags.Pass)
			cache.SaveMetaSession(address)
		} else {
			metas := strings.Split(loginFlags.Metas, ",")
			if len(metas) == 0 {
				return fmt.Errorf("meta server address is empty")
			}
			_, err := meta.NewMetaClient(loginFlags.Metas)
			if err != nil {
				return err
			}
			cache.SaveMetaSession(loginFlags.Metas)
		}
		return nil
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
	Use:   "createcluster",
	Short: "Create cluster in meta server.",
	Long:  `nebula-console createcluster --cluster [clustername] --replica [replica] --zones [zone1,zone2,...] --if_not_exists`,
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
		if resp.GetErrorCode() != 0 {
			return metaConsoleError("Create cluster failed", resp.GetErrorMsg())
		}
		fmt.Println("Create cluster successfully.")
		return nil
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
	Use:   "addservice",
	Short: `Add service into assigned cluster.`,
	Long:  `nebula-console meta addservice --type [graph|storage] --ip [ip] --port [port] --cluster [clustername]`,
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
		if resp.GetErrorCode() != 0 {
			return metaConsoleError("Add service failed", resp.GetErrorMsg())
		}
		fmt.Println("Add service successfully.")
		return nil
	},
}

type InitClusterFlags struct {
	Cluster string
}

var initClusterFlags InitClusterFlags

var metaInitClusterCmd = &cobra.Command{
	Use:   "initcluster",
	Short: "Init cluster storage part.",
	Long:  `nebula-console meta initcluster --cluster [clustername]`,
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
		if resp.Code != 0 {
			return metaConsoleError("Init cluster failed", resp.GetErrorMsg())
		}
		fmt.Println("Init cluster successfully.")
		return nil
	},
}

type ShowServiceFlags struct {
	Cluster string
}

var showServiceFlags ShowServiceFlags

var metaShowServiceCmd = &cobra.Command{
	Use:   "showservice",
	Short: "Show service in cluster.",
	Long:  `nebula-console meta showservice --cluster [clustername]`,
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
		if resp.Code != 0 {
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

type ShowClusterFlags struct {
	Cluster string
}

var showClusterFlags ShowClusterFlags

var metaShowClusterCmd = &cobra.Command{
	Use:   "showcluster",
	Short: "Show cluster, show all if no cluster name specified.",
	Long:  `nebula-console meta showcluster --cluster [clustername]`,
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
		if resp.GetErrorCode() != 0 {
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

func init() {
	metaCmd.AddCommand(metaLoginCmd)
	metaCmd.AddCommand(metaCreateClusterCmd)
	metaCmd.AddCommand(metaAddServiceCmd)
	metaCmd.AddCommand(metaInitClusterCmd)
	metaCmd.AddCommand(metaShowServiceCmd)
	metaCmd.AddCommand(metaShowClusterCmd)

	metaLoginCmd.Flags().StringVarP(&loginFlags.Addr, "addr", "a", "", "meta server address")
	metaLoginCmd.Flags().Uint32VarP(&loginFlags.Port, "port", "", 0, "meta server port")
	metaLoginCmd.Flags().StringVarP(&loginFlags.User, "user", "u", "", "user name")
	metaLoginCmd.Flags().StringVarP(&loginFlags.Pass, "password", "p", "", "password")
	metaLoginCmd.Flags().StringVarP(&loginFlags.Metas, "metas", "m", "", "meta server address list separated by comma like \"xx:xx,xx:xx\"")

	metaCreateClusterCmd.Flags().StringVarP(&createClusterFlags.ClusterName, "cluster", "c", "", "cluster name")
	metaCreateClusterCmd.Flags().IntVarP(&createClusterFlags.Replica, "replica", "r", 0, "replica number")
	metaCreateClusterCmd.Flags().StringArrayVarP(&createClusterFlags.Zones, "zones", "z", []string{}, "zones")
	metaCreateClusterCmd.Flags().BoolVarP(&createClusterFlags.IfNotExists, "if_not_exists", "", false, "if not exists")

	metaAddServiceCmd.Flags().StringVarP(&addServiceFlags.ServiceType, "type", "t", "", "service type")
	metaAddServiceCmd.Flags().StringVarP(&addServiceFlags.IP, "ip", "i", "", "service ip")
	metaAddServiceCmd.Flags().Uint32VarP(&addServiceFlags.Port, "port", "p", 0, "service port")
	metaAddServiceCmd.Flags().StringVarP(&addServiceFlags.Cluster, "cluster", "c", "", "cluster name")

	metaInitClusterCmd.Flags().StringVarP(&initClusterFlags.Cluster, "cluster", "c", "", "cluster name")
	metaShowServiceCmd.Flags().StringVarP(&showServiceFlags.Cluster, "cluster", "c", "", "cluster name")
	metaShowClusterCmd.Flags().StringVarP(&showClusterFlags.Cluster, "cluster", "c", "", "cluster name")
}

func main() {
	metaCmd.Execute()
}
