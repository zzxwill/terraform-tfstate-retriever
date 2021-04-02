package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	Namespace             = "CONFIGMAPS_NAMESPACE"
	TFStateConfigMapsName = "CONFIGMAPS_NAME"
	TFStateDir            = "TF_STATE_DIR"
	TFStateName           = "TF_STATE_NAME"
)

const TerraformStateName = "terraform.tfstate"

func main() {
	var (
		namespace             = os.Getenv(Namespace)
		tfStateConfigMapsName = os.Getenv(TFStateConfigMapsName)
		tfStateDir            = os.Getenv(TFStateDir)
		tfStateName           = os.Getenv(TFStateName)
	)
	if namespace == "" {
		namespace = "default"
	}

	if tfStateConfigMapsName == "" {
		tfStateConfigMapsName = "poc-tf-state"
	}

	if tfStateDir == "" {
		tfStateDir = "/data/tfstate"
	}

	if tfStateName == "" {
		tfStateName = TerraformStateName
	}
	ctx := context.Background()

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var tfStateCM = v1.ConfigMap{
		TypeMeta:   metav1.TypeMeta{APIVersion: "v1", Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{Name: tfStateConfigMapsName, Namespace: namespace},
	}

	tfStateFile := filepath.Join(tfStateDir, tfStateName)

	for {
		time.Sleep(5)
		fmt.Printf("checking file %s\n", tfStateFile)
		state, err := ioutil.ReadFile(tfStateFile)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		tfStateCM.Data = map[string]string{TerraformStateName: string(state)}
		existedCM, err := clientSet.CoreV1().ConfigMaps(namespace).Get(ctx, tfStateConfigMapsName, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				fmt.Printf("ConfigMaps %s doesn't exist, trying to create it\n", tfStateConfigMapsName)
				if _, err := clientSet.CoreV1().ConfigMaps(namespace).Create(ctx, &tfStateCM, metav1.CreateOptions{}); err != nil {
					fmt.Println(err.Error())
					continue
				}
				break
			}
			fmt.Println(err.Error())
			continue
		}
		existedCM.Data = map[string]string{TerraformStateName: string(state)}
		fmt.Printf("ConfigMaps %s exists, trying to update it\n", tfStateConfigMapsName)
		if _, err := clientSet.CoreV1().ConfigMaps(namespace).Update(ctx, &tfStateCM, metav1.UpdateOptions{}); err != nil {
			fmt.Println(err.Error())
			continue
		}
		break
	}

}
