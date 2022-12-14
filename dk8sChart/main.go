package main

import (
	"example.com/cdk8s/imports/k8s"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
)

const APP_NAME = "torrent-z"
const INGRESS_PATH = "z.torrent.zxyxyhome.duckdns.org"
const NAMESPACE = "torrent-z"
const NFS_PATH = "/storage02/kube/torrent-z"
const NFS_SERVER = "172.21.0.2"

type MyChartProps struct {
	cdk8s.ChartProps
}

func NewMyChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	label := map[string]*string{"app": jsii.String(APP_NAME)}

	udp := "UDP"
	// define resources here
	service := k8s.NewKubeService(chart, jsii.String("service"), &k8s.KubeServiceProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String("transmission"),
			Namespace: jsii.String(NAMESPACE),
		},
		Spec: &k8s.ServiceSpec{
			Type: jsii.String("ClusterIP"),
			Ports: &[]*k8s.ServicePort{{
				Name:       jsii.String("web"),
				Port:       jsii.Number(80),
				TargetPort: k8s.IntOrString_FromNumber(jsii.Number(80)),
			}, {
				Name:       jsii.String("upload"),
				Port:       jsii.Number(20001),
				TargetPort: k8s.IntOrString_FromNumber(jsii.Number(20001)),
				Protocol:   &udp,
			}},
			Selector: &label,
		},
	})

	k8s.NewKubeDeployment(chart, jsii.String("deployment"), &k8s.KubeDeploymentProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String("transmission"),
			Namespace: jsii.String(NAMESPACE),
		},
		Spec: &k8s.DeploymentSpec{
			Replicas: jsii.Number(1),
			Selector: &k8s.LabelSelector{
				MatchLabels: &label,
			},
			Template: &k8s.PodTemplateSpec{
				Metadata: &k8s.ObjectMeta{
					Labels: &label,
				}, Spec: &k8s.PodSpec{
					Containers: &[]*k8s.Container{{
						Name:  jsii.String("transmission"),
						Image: jsii.String("zxyxy/transmission:0.2.0"),
						Ports: &[]*k8s.ContainerPort{{
							ContainerPort: jsii.Number(80),
						}, {
							ContainerPort: jsii.Number(20001),
							Protocol:      &udp,
						}},
						VolumeMounts: &[]*k8s.VolumeMount{{
							Name:      jsii.String("settingsjson"),
							MountPath: jsii.String("/etc/transmission/settings.json"),
							SubPath:   jsii.String("settings.json"),
						}},
						Env: &[]*k8s.EnvVar{{
							Name:  jsii.String("TRANS_UID"),
							Value: jsii.String("1000"),
						}, {
							Name:  jsii.String("TRANS_GID"),
							Value: jsii.String("1000"),
						}},
					}},
					Volumes: &[]*k8s.Volume{{
						Name: jsii.String("settingsjson"),
						ConfigMap: &k8s.ConfigMapVolumeSource{
							Name: jsii.String("transmission"),
						},
					}},
				},
			},
		},
	})

	k8s.NewKubeIngress(chart, jsii.String("ingress"), &k8s.KubeIngressProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String("transmission"),
			Namespace: jsii.String(NAMESPACE),
		},
		Spec: &k8s.IngressSpec{
			Rules: &[]*k8s.IngressRule{{
				Host: jsii.String(INGRESS_PATH),
				Http: &k8s.HttpIngressRuleValue{
					Paths: &[]*k8s.HttpIngressPath{{
						PathType: jsii.String("Prefix"),
						Path:     jsii.String("/"),
						Backend: &k8s.IngressBackend{
							Service: &k8s.IngressServiceBackend{
								Name: jsii.String(*service.Name()),
								Port: &k8s.ServiceBackendPort{
									Number: jsii.Number(80),
								},
							},
						}},
					},
				},
			}},
		},
	})

	k8s.NewKubePersistentVolume(chart, jsii.String("torrent-pv"), &k8s.KubePersistentVolumeProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String("torrent"),
		},
		Spec: &k8s.PersistentVolumeSpec{
			Capacity: &map[string]k8s.Quantity{
				"storage": k8s.Quantity_FromString(jsii.String("2500Gi")),
			},
			VolumeMode:                    jsii.String("Filesystem"),
			PersistentVolumeReclaimPolicy: jsii.String("Retain"),
			StorageClassName:              jsii.String("torrent-z"),
			MountOptions:                  jsii.Strings("hard", "nfsvers=4.1"),
			Nfs: &k8s.NfsVolumeSource{
				Path:   jsii.String(NFS_PATH),
				Server: jsii.String(NFS_SERVER),
			},
			AccessModes: jsii.Strings("ReadWriteOnce"),
		},
	})

	k8s.NewKubePersistentVolumeClaim(chart, jsii.String("torrent-pvc"), &k8s.KubePersistentVolumeClaimProps{
		Metadata: &k8s.ObjectMeta{
			Name:      jsii.String("torrent-pvc"),
			Namespace: jsii.String(NAMESPACE),
		},
		Spec: &k8s.PersistentVolumeClaimSpec{
			VolumeMode:       jsii.String("Filesystem"),
			StorageClassName: jsii.String("torrent-z"),
			Resources: &k8s.ResourceRequirements{
				Requests: &map[string]k8s.Quantity{
					"storage": k8s.Quantity_FromString(jsii.String("2500Gi")),
				},
			},
			AccessModes: jsii.Strings("ReadWriteOnce"),
		},
	})

	k8s.NewKubeNamespace(chart, jsii.String("namespace"), &k8s.KubeNamespaceProps{
		Metadata: &k8s.ObjectMeta{
			Name: jsii.String(NAMESPACE),
		},
		Spec: &k8s.NamespaceSpec{},
	})
	return chart
}

func main() {
	app := cdk8s.NewApp(nil)
	NewMyChart(app, "torrent-z", nil)
	app.Synth()
}
