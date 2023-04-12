package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kubeedge/examples/kubeedge-counter-demo/web-controller-app/utils"

	devices "github.com/kubeedge/kubeedge/pkg/apis/devices/v1alpha2"

	"github.com/astaxie/beego"

	"k8s.io/client-go/rest"
)

type DeviceStatus struct {
	Status devices.DeviceStatus `json:"status"`
}

// The device id of the counter
var deviceID = "yolocloud"

// The default namespace in which the counter device instance resides
var namespace = "default"

// The default status of the counter
var originCmd = "STATUS"

// The CRD client used to patch the device instance.
// rest.RESTClient 是 Kubernetes 提供的 REST 客户端库，用于与 Kubernetes API Server 进行交互
var crdClient *rest.RESTClient
var arr = make([]string, 0)
var num []string
var status = map[string]string{ //status 变量用于存储解析后的设备状态信息。
	//"status": "OFF",
	"imageid":   "0",
	"masknum":   "0",
	"nomasknum": "0",
}

// func init() {
// 	// Create a client to talk to the K8S API server to patch the device CRDs
// 	kubeConfig, err := utils.KubeConfig()
// 	if err != nil {
// 		log.Fatalf("Failed to create KubeConfig, error : %v", err)
// 	}
// 	log.Println("Get kubeConfig successfully")

// 	crdClient, err = utils.NewCRDClient(kubeConfig)
// 	if err != nil {
// 		log.Fatalf("Failed to create device crd client , error : %v", err)
// 	}
// 	log.Println("Get crdClient successfully")
// }

func UpdateStatus() map[string]string {
	result := DeviceStatus{}
	raw, _ := crdClient.Get().Namespace(namespace).Resource(utils.ResourceTypeDevices).Name(deviceID).DoRaw(context.TODO())
	_ = json.Unmarshal(raw, &result)
	for _, twin := range result.Status.Twins {
		// status["status"] = twin.Desired.Value
		// 使用 split() 方法将字符串拆分成一个列表
		num = strings.Split(twin.Reported.Value, "")
		status["imageid"] = num[0]
		status["masknum"] = num[1]
		status["nomasknum"] = num[2]
	}

	return status
}

type TrackController struct {
	beego.Controller
}

func (controller *TrackController) Status() {
	// 生成表格的 HTML 代码
	var output strings.Builder
	output.WriteString("<table style=\"border-collapse: collapse; width: 100%; text-align: center;\">\n")
	output.WriteString("<thead>\n<tr style=\"background-color: #eee;\">\n<th>ID</th>\n<th>Image ID</th>\n<th>Mask Num</th>\n<th>No Mask Num</th>\n</tr>\n</thead>\n")
	output.WriteString("<tbody>\n")
	result := DeviceStatus{}
	raw, _ := crdClient.Get().Namespace(namespace).Resource(utils.ResourceTypeDevices).Name(deviceID).DoRaw(context.TODO())
	_ = json.Unmarshal(raw, &result)
	for _, twin := range result.Status.Twins {
		arr = append(arr, twin.Reported.Value)
	}

	for i, str := range arr {
		items := strings.Split(str, ",")
		output.WriteString(fmt.Sprintf("<tr style=\"background-color: %s;\">\n", "#fff"))
		output.WriteString(fmt.Sprintf("<td style=\"border: 1px solid #ccc; padding: 5px;\">%d</td>\n", i+1))
		for _, v := range items {
			output.WriteString(fmt.Sprintf("<td style=\"border: 1px solid #ccc; padding: 5px;\">%s</td>\n", v))
		}
		output.WriteString("</tr>\n")
	}
	output.WriteString("</tbody>\n")
	output.WriteString("</table>\n")

	// 将生成的 HTML 代码输出到浏览器
	controller.Ctx.WriteString(output.String())
}

// Index is the initial view
func (controller *TrackController) Index() {
	log.Println("Index Start")

	controller.Layout = "layout.html"
	controller.TplName = "content.html"
	controller.LayoutSections = map[string]string{}
	controller.LayoutSections["PageHead"] = "head.html"
	log.Println("Index Finish")
}

// Control
func (controller *TrackController) ControlTrack() {
	// Get track id
	params := struct {
		TrackID string `form:":trackId"`
	}{controller.GetString(":trackId")}

	resultCode := 0

	status := map[string]string{}

	log.Printf("ControlTrack: %s", params.TrackID)

	if params.TrackID == "STATUS" {
		status = UpdateStatus()
	}

	// response
	controller.AjaxResponse(resultCode, status, nil)
}

func (Controller *TrackController) AjaxResponse(resultCode int, resultString map[string]string, data interface{}) {
	response := struct {
		Result       int
		ResultString map[string]string
		ResultObject interface{}
	}{
		Result:       resultCode,
		ResultString: resultString,
		ResultObject: data,
	}

	Controller.Data["json"] = response
	Controller.ServeJSON()
}
