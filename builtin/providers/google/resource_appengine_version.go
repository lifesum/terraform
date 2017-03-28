package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	appengine "google.golang.org/api/appengine/v1"
)

func resourceAppEngineVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppEngineVersionCreate,
		Read:   resourceAppEngineVersionRead,
		Update: resourceAppEngineVersionUpdate,
		Delete: resourceAppEngineVersionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// ProjectId: [Optional] The ID of the project containing this dataset.
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"runtime": {
				Type:     schema.TypeString,
				Required: true,
			},

			"apps_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"services_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},

			"automatic_scaling": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{"basic_scaling"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// CoolDownPeriod: Amount of time that the Autoscaler
						// (https://cloud.google.com/compute/docs/autoscaler/) should wait
						// between changes to the number of virtual machines. Only applicable
						// for VM runtimes.
						"cooldown_period": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// CpuUtilization: Target scaling by CPU usage.
						"cpu_utilization": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// AggregationWindowLength: Period of time over which CPU utilization is
									// calculated.
									"aggregation_window_length": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"target_utilization": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						// DiskUtilization: Target scaling by disk usage.
						"disk_utilization": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// TargetReadBytesPerSecond: Target bytes read per second.
									"target_read_bytes_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetReadOpsPerSecond: Target ops read per seconds.
									"target_read_ops_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetWriteBytesPerSecond: Target bytes written per second.
									"target_write_bytes_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetWriteOpsPerSecond: Target ops written per second.
									"target_write_ops_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						// MaxConcurrentRequests: Number of concurrent requests an automatic
						// scaling instance can accept before the scheduler spawns a new
						// instance.Defaults to a runtime-specific value.
						"max_concurrent_requests": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// MaxIdleInstances: Maximum number of idle instances that should be
						// maintained for this version.
						"max_idle_instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// MaxPendingLatency: Maximum amount of time that a request should wait
						// in the pending queue before starting a new instance to handle it.
						"max_pending_latency": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Automatic",
						},

						// MaxTotalInstances: Maximum number of instances that should be started
						// to handle requests.
						"max_total_instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// MinIdleInstances: Minimum number of idle instances that should be
						// maintained for this version. Only applicable for the default version
						// of a service.
						"min_idle_instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// MinPendingLatency: Minimum amount of time a request should wait in
						// the pending queue before starting a new instance to handle it.
						"min_pending_latency": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Automatic",
						},

						// MinTotalInstances: Minimum number of instances that should be
						// maintained for this version.
						"min_total_instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// NetworkUtilization: Target scaling by network usage.
						"network_utilization": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// TargetReceivedBytesPerSecond: Target bytes received per second.
									"target_received_bytes_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetReceivedPacketsPerSecond: Target packets received per second.
									"target_received_packets_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetSentBytesPerSecond: Target bytes sent per second.
									"target_sent_bytes_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetSentPacketsPerSecond: Target packets sent per second.
									"target_sent_packages_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},

						// RequestUtilization: Target scaling by request utilization.
						"request_utilization": {
							Type:     schema.TypeList,
							MinItems: 1,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// TargetConcurrentRequests: Target number of concurrent requests.
									"target_concurrent_requests": {
										Type:     schema.TypeInt,
										Optional: true,
									},

									// TargetRequestCountPerSecond: Target requests per second.
									"target_request_count_per_second": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"basic_scaling": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MinItems:      1,
				MaxItems:      1,
				ConflictsWith: []string{"automatic_scaling"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// IdleTimeout: Duration of time after the last request that an instance
						// must wait before the instance is shut down.
						"idle_timeout": {
							Type:     schema.TypeString,
							Optional: true,
						},

						// MaxInstances: Maximum number of instances to create for this version.
						"max_instances": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			// BetaSettings: Metadata settings that are supplied to this version to
			// enable beta runtime features.
			"beta_settings": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},

			// CreateTime: Time that this version was created.@OutputOnly
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// CreatedBy: Email address of the user who created this
			// version.@OutputOnly
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// InboundServices: Before an application can receive email or XMPP
			// messages, the application must be configured to enable the service.
			//
			// Possible values:
			//   "INBOUND_SERVICE_UNSPECIFIED" - Not specified.
			//   "INBOUND_SERVICE_MAIL" - Allows an application to receive mail.
			//   "INBOUND_SERVICE_MAIL_BOUNCE" - Allows an application to receive
			// email-bound notifications.
			//   "INBOUND_SERVICE_XMPP_ERROR" - Allows an application to receive
			// error stanzas.
			//   "INBOUND_SERVICE_XMPP_MESSAGE" - Allows an application to receive
			// instant messages.
			//   "INBOUND_SERVICE_XMPP_SUBSCRIBE" - Allows an application to receive
			// user subscription POSTs.
			//   "INBOUND_SERVICE_XMPP_PRESENCE" - Allows an application to receive
			// a user's chat presence.
			//   "INBOUND_SERVICE_CHANNEL_PRESENCE" - Registers an application for
			// notifications when a client connects or disconnects from a channel.
			//   "INBOUND_SERVICE_WARMUP" - Enables warmup requests.
			"inbound_services": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Threadsafe: Whether multiple requests can be dispatched to this
			// version at once.
			"threadsafe": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			// Deployment: Code and application artifacts that make up this
			// version.Only returned in GET requests if view=FULL is set.
			"deployment": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Container: The Docker image for the container that runs the version.
						// Only applicable for instances running in the App Engine flexible
						// environment.
						"container": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Image: URI to the hosted container image in Google Container
									// Registry. The URI must be fully qualified and include a tag or
									// digest. Examples: "gcr.io/my-project/image:tag" or
									// "gcr.io/my-project/image@digest"
									"image": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func expandDeployment(configured interface{}) *appengine.Deployment {
	raw := configured.([]interface{})[0].(map[string]interface{})
	deploy := &appengine.Deployment{}

	bucketName := "lifesum-terraform.appspot.com"
	bucketKey := "index.html"
	bucketURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, bucketKey)

	files := make(map[string]appengine.FileInfo)
	files[bucketURL] = appengine.FileInfo{
		SourceUrl: bucketURL,
	}
	deploy.Files = files

	if v, ok := raw["container"]; ok {
		container := v.([]interface{})[0].(map[string]interface{})
		deploy.Container.Image = container["image"].(string)
	}

	return deploy
}

func expandBasicScaling(configured interface{}) *appengine.BasicScaling {
	raw := configured.([]interface{})[0].(map[string]interface{})
	scaling := &appengine.BasicScaling{}

	if v, ok := raw["idle_timeout"]; ok {
		scaling.IdleTimeout = v.(string)
	}

	if v, ok := raw["max_instances"]; ok {
		scaling.MaxInstances = int64(v.(int))
	}

	return scaling
}

func expandAutomaticScaling(configured interface{}) *appengine.AutomaticScaling {
	raw := configured.([]interface{})[0].(map[string]interface{})
	scaling := &appengine.AutomaticScaling{}

	if v, ok := raw["cooldown_period"]; ok {
		scaling.CoolDownPeriod = v.(string)
	}

	if v, ok := raw["min_idle_instances"]; ok {
		scaling.MinIdleInstances = int64(v.(int))
	}

	if v, ok := raw["max_idle_instances"]; ok {
		scaling.MaxIdleInstances = int64(v.(int))
	}

	if v, ok := raw["min_pending_latency"]; ok {
		scaling.MinPendingLatency = v.(string)
	}

	if v, ok := raw["max_pending_latency"]; ok {
		scaling.MaxPendingLatency = v.(string)
	}

	return scaling
}

func resourceAppEngineVersionR(d *schema.ResourceData, meta interface{}) *appengine.Version {
	bucketName := "lifesum-terraform.appspot.com"
	bucketKey := "index.html"
	bucketURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, bucketKey)

	version := &appengine.Version{}

	return &appengine.Version{
		BasicScaling:    &appengine.BasicScaling{MaxInstances: 3},
		Deployment:      &appengine.Deployment{Files: map[string]appengine.FileInfo{"index.html": appengine.FileInfo{SourceUrl: bucketURL}}},
		Id:              "foobaaz",
		Runtime:         "python27",
		InboundServices: []string{"INBOUND_SERVICE_WARMUP"},
		EnvVariables:    map[string]string{"foo": "true"},
		Threadsafe:      false,
	}

	//
	// if v, ok := d.GetOk("runtime"); ok {
	// 	appVersion.Runtime = v.(string)
	// }
	//
	// if v, ok := d.GetOk("deployment"); ok {
	// 	appVersion.Deployment = expandDeployment(v)
	// }
	//
	// if v, ok := d.GetOk("threadsafe"); ok {
	// 	appVersion.Threadsafe = v.(bool)
	// }
	//
	// if v, ok := d.GetOk("inbound_services"); ok {
	// 	var services []string
	// 	for _, service := range v.([]interface{}) {
	// 		services = append(services, service.(string))
	// 	}
	// 	appVersion.InboundServices = services
	// }
	//
	// if v, ok := d.GetOk("automatic_scaling"); ok {
	// 	appVersion.AutomaticScaling = expandAutomaticScaling(v)
	// }
	//
	// if v, ok := d.GetOk("basic_scaling"); ok {
	// 	appVersion.BasicScaling = expandBasicScaling(v)
	// }
	// return appVersion
}

func resourceAppEngineVersionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	version := resourceAppEngineVersionR(d, meta)

	j, _ := version.MarshalJSON()
	return fmt.Errorf("JSON: \n%s", j)
	appsID := d.Get("apps_id").(string)

	op, err := config.clientAppEngine.Apps.Services.Versions.Create(appsID+"/"+"default", "default", version).Do()
	if err != nil {
		return err
	}

	waitErr := appEngineOperationWait(config, op, appsID, "creating AppEngine service", 30, 3)
	if waitErr != nil {
		d.SetId("")
		return waitErr
	}

	return resourceAppEngineVersionRead(d, meta)
}

func resourceAppEngineVersionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	appsID := d.Get("apps_id").(string)
	servicesID := d.Get("services_id").(string)
	versionsID := d.Get("version").(string)

	res, err := config.clientAppEngine.Apps.Services.Versions.Get(appsID, servicesID, versionsID).Do()
	if err != nil {
		return err
	}

	d.SetId(res.Id)

	return nil
}

func resourceAppEngineVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceAppEngineVersionDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
