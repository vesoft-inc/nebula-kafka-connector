package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vesoft-inc/nebula-ng-tools/console/meta"
)

var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "Execute meta command in cli mode.",
	Long: `Execute meta command in cli mode. Use 'nebula-console meta -h' to see usage.
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

var metaLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login meta server.",
	Long:  `login meta server --addr [ip] --port [port] --user [user] --password [password]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loginFlags.Metas == "" {
			address := fmt.Sprintf("%s:%d", loginFlags.Addr, loginFlags.Port)
			metaclient := meta.NewMetaClient(address)
			leader, err := metaclient.GetLeader(meta.NewGetLeaderRequest())
			if err != nil {
				return fmt.Errorf("get leader failed: %s", err)
			}
			metaclient = meta.NewMetaClient(fmt.Sprintf("%s:%d", leader.Host, leader.Port))
			metaclient.Login(loginFlags.User, loginFlags.Pass)
			meta.SaveMetaSession(address, fmt.Sprintf("%s:%d", leader.Host, leader.Port))
			return nil
		} else {
			metas := strings.Split(loginFlags.Metas, ",")
			if len(metas) == 0 {
				return fmt.Errorf("meta server address is empty")
			}
			for _, m := range metas {
				metaclient := meta.NewMetaClient(m)
				leader, err := metaclient.GetLeader(meta.NewGetLeaderRequest())
				if err != nil {
					fmt.Println("get leader failed: ", err)
					continue
				}
				metaclient = meta.NewMetaClient(fmt.Sprintf("%s:%d", leader.Host, leader.Port))
				metaclient.Login(loginFlags.User, loginFlags.Pass)
				meta.SaveMetaSession(metas[0], fmt.Sprintf("%s:%d", leader.Host, leader.Port))
				return nil
			}

		}
		return fmt.Errorf("get leader failed.")
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
	RunE: func(cmd *cobra.Command, args []string) error {
		metaclient, err := meta.LoadMetaClient()
		if err != nil {
			return fmt.Errorf("load meta session failed: %s", err)
		}
		defer metaclient.Close()

		flags := createClusterFlags
		if flags.ClusterName == "" {
			return fmt.Errorf("cluster name is empty")
		}
		if flags.Replica == 0 {
			return fmt.Errorf("replica number is invalid")
		}
		req := meta.NewCreateClusterRequest(flags.ClusterName, flags.Replica, flags.Zones, flags.IfNotExists)

		msg, _, err := metaclient.CreateCluster(req)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("create cluster failed")
		}
		fmt.Println(msg)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		metaclient, err := meta.LoadMetaClient()
		if err != nil {
			return fmt.Errorf("load meta session failed: %s", err)
		}
		defer metaclient.Close()

		flags := addServiceFlags

		var serviceType int8
		if flags.ServiceType == "graph" {
			serviceType = meta.Graph
		} else if flags.ServiceType == "storage" {
			serviceType = meta.Storage
		} else {
			return fmt.Errorf("service type is not correct")
		}

		if flags.IP == "" {
			return fmt.Errorf("service ip is empty")
		}

		if flags.Cluster == "" {
			return fmt.Errorf("cluster name is empty")
		}
		req := meta.NewAddServiceRequest(flags.IP, flags.Port, serviceType, flags.Cluster)

		msg, _, err := metaclient.AddService(req)
		if err != nil {
			return fmt.Errorf("add service failed")
		}
		fmt.Println(msg)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		metaclient, err := meta.LoadMetaClient()
		if err != nil {
			return fmt.Errorf("load meta session failed: %s", err)
		}
		defer metaclient.Close()

		cluster := initClusterFlags.Cluster
		if cluster == "" {
			return fmt.Errorf("cluster name is empty")
		}
		req := meta.NewInitClusterRequest(cluster)

		msg, _, err := metaclient.InitCluster(req)
		if err != nil {
			return fmt.Errorf("init cluster failed: %s", err)
		}
		fmt.Println(msg)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		metaclient, err := meta.LoadMetaClient()
		if err != nil {
			return fmt.Errorf("load meta session failed: %s", err)
		}
		defer metaclient.Close()

		cluster := showServiceFlags.Cluster
		if cluster == "" {
			return fmt.Errorf("cluster name is empty")
		}
		req := meta.NewListServiceRequest(cluster)

		msg, _, err := metaclient.ShowService(req)
		if err != nil {
			return fmt.Errorf("list service failed: %s", err)
		}
		fmt.Println(msg)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		metaclient, err := meta.LoadMetaClient()
		if err != nil {
			return fmt.Errorf("load meta session failed: %s", err)
		}
		defer metaclient.Close()

		cluster := showClusterFlags.Cluster
		req := meta.NewListClusterRequest(cluster)

		msg, _, err := metaclient.ShowCluster(req)
		if err != nil {
			return fmt.Errorf("list cluster failed: %s", err)
		}
		fmt.Println(msg)
		return nil
	},
}

type CreateSchemaFlags struct {
	Name        string
	Path        string
	IfNotExists bool
}

var createSchemaFlags CreateSchemaFlags

var metaCreateSchemaCmd = &cobra.Command{
	Use:   "createschema",
	Short: "Create schema in catalog tree.",
	Long:  `nebula-console meta createschema --name [name] --path [path] --if_not_exists`,
	RunE: func(cmd *cobra.Command, args []string) error {
		metaclient, err := meta.LoadMetaClient()
		if err != nil {
			return fmt.Errorf("load meta session failed: %s", err)
		}
		defer metaclient.Close()

		flags := createSchemaFlags

		if flags.Name == "" {
			return fmt.Errorf("schema name is empty")
		}

		if flags.Path == "" {
			return fmt.Errorf("schema path is empty")
		}

		// ifNotExists, _ := cmd.Flags().GetBool("if_not_exists")

		m := make(map[string]interface{})
		m["type"] = "createSchema"
		paths := []string{"/"}
		names := strings.Split(flags.Path, "/")[1:]
		for _, n := range names {
			if n != "" {
				paths = append(paths, n)
			}
		}
		m["path"] = paths
		m["name"] = flags.Name

		jsonBytes, err := json.Marshal(m)
		if err != nil {
			return fmt.Errorf("failed to format command: %s", err)
		}
		req := meta.NewMetaDDLRequest("createSchema", string(jsonBytes))

		msg, _, err := metaclient.CreateSchema(req)
		if err != nil {
			return fmt.Errorf("create schema failed: %s", err)
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)

	metaCmd.AddCommand(metaLoginCmd)
	metaCmd.AddCommand(metaCreateClusterCmd)
	metaCmd.AddCommand(metaAddServiceCmd)
	metaCmd.AddCommand(metaInitClusterCmd)
	metaCmd.AddCommand(metaShowServiceCmd)
	metaCmd.AddCommand(metaShowClusterCmd)
	metaCmd.AddCommand(metaCreateSchemaCmd)

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

	metaCreateSchemaCmd.Flags().StringVarP(&createSchemaFlags.Name, "name", "n", "", "schema name")
	metaCreateSchemaCmd.Flags().StringVarP(&createSchemaFlags.Path, "path", "p", "", "schema path")
	metaCreateSchemaCmd.Flags().BoolVarP(&createClusterFlags.IfNotExists, "if_not_exists", "", false, "if not exists")
}
