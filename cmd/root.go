package cmd

import (
	"fmt"
	"github.com/joostvdg/cat/pkg/api/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const (
	// apiVersion is version of API is provided by server
	apiVersion = "v1"
)

var ServerPort string
var ServerHost string
var ServerProtocol string
var Uuid string
var Name string
var Namespace string
var Description string
var Sources string
var ArtifactIDs string
var Labels string
var Annotations string

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)

	createCommand.Flags().StringVarP(&Name, "name", "n", "", "Name of the application")
	createCommand.Flags().StringVarP(&Namespace, "namespace", "c", "", "Namespace of the application")
	createCommand.Flags().StringVarP(&Description, "description", "d", "", "Description of the application")
	createCommand.Flags().StringVarP(&Sources, "sources", "s", "", "Sources of the application to perform the action on, ',' separated")
	createCommand.Flags().StringVarP(&ArtifactIDs, "artifactIDs", "r", "", "ArtifactIDs of the application',' separated")
	createCommand.Flags().StringVarP(&Labels, "labels", "l", "", "Labels of the application, ',' separated for the list, ';' separated for the key value")
	createCommand.Flags().StringVarP(&Annotations, "annotations", "a", "", "Annotations of the application, ',' separated for the list, ';' separated for the origin/key/value")

	getCommand.Flags().StringVarP(&Uuid, "uuid", "u", "", "Uuid of the application to perform the action on")

	rootCmd.Flags().StringVarP(&ServerPort, "port", "P", "9090", "Port to run the gRPC API server on")
	rootCmd.Flags().StringVarP(&ServerHost, "host", "H", "localhost", "Host to call the gRPC API server on")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(getCommand)
	rootCmd.AddCommand(getAllCommand)
	rootCmd.AddCommand(deleteCommand)
	rootCmd.AddCommand(createCommand)
}

var getCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieves a single entry",
	Long:  `This will call the gRPC API server of CAT, and retrieve a single Application entry`,
	Run: func(cmd *cobra.Command, args []string) {
		address := ServerHost + ":" + ServerPort
		GetApplication(address, Uuid)
	},
}

var createCommand = &cobra.Command{
	Use:   "create",
	Short: "Creates a new entry",
	Long:  `This will call the gRPC API server of CAT, and create a single Application entry`,
	Run: func(cmd *cobra.Command, args []string) {
		address := ServerHost + ":" + ServerPort

		// strings.Split("127.0.0.1:5432", ":")
		sources := strings.Split(Sources, ",")
		artifactIDs := strings.Split(ArtifactIDs, ",")
		labelPairs := strings.Split(Labels, ",")
		labels := make([]*v1.Label, 0)
		for _, labelPair := range labelPairs {
			if strings.Contains(labelPair, ";") {
				keyValue := strings.Split(labelPair, ";")
				label := v1.Label{
					Key:   keyValue[0],
					Value: keyValue[1],
				}
				labels = append(labels, &label)
			}
		}

		annotationGroups := strings.Split(Annotations, ",")
		annotations := make([]*v1.Annotation, 0)
		for _, annotationGroup := range annotationGroups {
			if strings.Contains(annotationGroup, ";") {
				originKeyValue := strings.Split(annotationGroup, ";")
				annotation := v1.Annotation{
					Origin: originKeyValue[0],
					Key:    originKeyValue[1],
					Value:  originKeyValue[2],
				}
				annotations = append(annotations, &annotation)
			}
		}

		application := v1.Application{
			Name:        Name,
			Description: Description,
			Namespace:   Namespace,
			Sources:     sources,
			ArtifactIDs: artifactIDs,
			Labels:      labels,
			Annotations: annotations,
		}

		CreateApplication(address, application)
	},
}

var deleteCommand = &cobra.Command{
	Use:   "get",
	Short: "Retrieves a single entry",
	Long:  `This will call the gRPC API server of CAT, and retrieve a single Application entry`,
	Run: func(cmd *cobra.Command, args []string) {
		address := ServerHost + ":" + ServerPort
		DeleteApplication(address, Uuid)
	},
}

var getAllCommand = &cobra.Command{
	Use:   "all",
	Short: "Retrieves a all entries",
	Long:  `This will call the gRPC API server of CAT, and retrieve all Application entries`,
	Run: func(cmd *cobra.Command, args []string) {
		address := ServerHost + ":" + ServerPort
		GetAll(address)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CAT",
	Long:  `All software has versions. This is CAT's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CAT 0.1.0")
	},
}

var rootCmd = &cobra.Command{
	Use:   "cat-grpc-client",
	Short: "cat is a small application tracker",
	Long:  `Yada yada yada`,
	Run: func(cmd *cobra.Command, args []string) {
		// return "0.1.0"
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
