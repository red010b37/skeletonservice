package skelconfig

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/kms/apiv1"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

type MapConfig map[string]interface{}

var AppConfig appConfig
var ENV env
var ApiVersion = "1.0.0"

var DBconn *sql.DB

type env struct {
	ProjectId string

	HostName     string
	InstanceId   string
	DeploymentId string
	ServiceName  string
	Version      string
	IsLocal      bool
}

type appConfig struct {
	//DBHost     string
	//DBUser     string
	//DBPassword string
	//DBName     string
	//DBPort     int

	JWTSecret    string
	ConfigUrl    string
	ConfigApiKey string
}

var kmsClient *kms.KeyManagementClient

// init is auto run on import - this float64sets up the whole config for the app
// It also check for all the runtime vars and db connections required
// to make the auth service run
func init() {

	AppConfig = appConfig{}

	getEnvVars()

	// Check to see if system is running tests on travis

	if isRunningTests() {
		return
	}

	if ENV.IsLocal {
		log.Println("IsLocal IS set looking env service overrides")
		discoverLocalConfigService()
	} else {

		log.Println("IsLocal is NOT set connecting to config service")
		discoverConfigService()
	}

	loadConfigFromConfigService()

	errArr := checkConfig()
	if len(errArr) != 0 {

		for _, err := range errArr {
			log.Println(err.Error())
		}

		log.Fatal("Missing Config vars above")

	}

	//checkDBConnection()
}

func isRunningTests() bool {

	if viper.GetString("RACK_ENV") == "test" {
		return true
	} else {
		return false
	}

}

func discoverLocalConfigService() {

	configURL := viper.GetString("GLO_CONFIG_URL")

	if configURL == "" {
		log.Fatal("GLO_CONFIG_URL is not set")
	}

	AppConfig.ConfigUrl = configURL

}

func loadConfigFromConfigService() {

	url := fmt.Sprintf("%s/v1/subscriptions?key=%s", AppConfig.ConfigUrl, AppConfig.ConfigApiKey)

	resp, err := http.Get(url)
	checkErr(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	toMap := MapConfig{}

	err = json.Unmarshal(body, &toMap)
	checkErr(err)

	data := toMap["data"].(map[string]interface{})

	AppConfig.JWTSecret = data["GLO_V1_JWT_SECRET"].(string)

}

func FirestoreClient() *firestore.Client {

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, ENV.ProjectId)
	checkErr(err)

	return client
}

func discoverConfigService() {

	client := FirestoreClient()
	// Close client when done.
	defer client.Close()

	ctx := context.Background()
	snapshot, snapErr := client.Collection("system-env").Doc("config-location").Get(ctx)
	checkErr(snapErr)

	url, urlDataErr := snapshot.DataAt("url")
	checkErr(urlDataErr)

	apiKey, apiErr := snapshot.DataAt("apiKey")
	checkErr(apiErr)

	AppConfig.ConfigApiKey = apiKey.(string)
	AppConfig.ConfigUrl = url.(string)

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getEnvVars() {

	ENV = env{}
	viper.SetDefault("IS_LOCAL", false)
	viper.SetDefault("GOOGLE_CLOUD_PROJECT", "logmate-platform-prod")

	// Read teh AppEngine settings - set local dev first
	viper.SetDefault("GAE_DEPLOYMENT_ID", 0)
	viper.SetDefault("GAE_INSTANCE", "local-instacne")
	viper.SetDefault("GAE_SERVICE", "local-subservice")
	viper.SetDefault("GAE_VERSION", "local-version")
	viper.SetDefault("HOSTNAME", "local-hostname")

	// read the end
	viper.AutomaticEnv()

	ENV.ProjectId = viper.GetString("GOOGLE_CLOUD_PROJECT")
	ENV.DeploymentId = viper.GetString("GAE_DEPLOYMENT_ID")
	ENV.InstanceId = viper.GetString("GAE_INSTANCE")
	ENV.ServiceName = viper.GetString("GAE_SERVICE")
	ENV.Version = viper.GetString("GAE_VERSION")
	ENV.HostName = viper.GetString("HOSTNAME")
	ENV.IsLocal = viper.GetBool("IS_LOCAL")

}

//func checkDBConnection() {
//
//	mode := "require"
//
//	if ENV.IsLocal {
//		log.Println("Is in local mode disabling ssl")
//		mode = "disable"
//	}
//
//	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
//		AppConfig.DBHost,
//		AppConfig.DBUser,
//		AppConfig.DBPassword,
//		AppConfig.DBName,
//		AppConfig.DBPort,
//		mode)
//
//	var err error
//	DBconn, err = sql.Open("postgres", dbinfo)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = DBconn.Ping()
//	if err != nil {
//		DBconn.Close()
//		log.Fatal(err)
//	}
//
//	log.Println("DB ping successful")
//
//}

func checkConfig() []error {

	errorArr := []error{}
	//
	//if AppConfig.DBHost == "" {
	//	errorArr = append(errorArr, ConfDBHostErr)
	//}
	//
	//if AppConfig.DBUser == "" {
	//	errorArr = append(errorArr, ConfDBUserErr)
	//}
	//
	//if AppConfig.DBPassword == "" {
	//	errorArr = append(errorArr, ConfDBPasswordErr)
	//}
	//
	//if AppConfig.DBName == "" {
	//	log.Println()
	//	errorArr = append(errorArr, ConfDBNameErr)
	//}
	//
	//if AppConfig.DBPort == 0 {
	//	errorArr = append(errorArr, ConfDBPortErr)
	//}
	//
	//if AppConfig.HashPepper == "" {
	//	errorArr = append(errorArr, ConfHashPepperErr)
	//}

	if AppConfig.JWTSecret == "" {
		errorArr = append(errorArr, ConfJWTSecretErr)
	}

	//if AppConfig.ConfigUrl == "" {
	//	errorArr = append(errorArr, ConfConfigURLErr)
	//}

	return errorArr
}
