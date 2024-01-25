package us1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateUs1MongoIntegrationDebuggerTemplate() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEbsDlmLifecyclePolicy"] = mongo.GetDlmLifeCyclePolicy()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("us1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("us1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "13.176/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.13.176 - 10.16.13.191
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "13.192/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.13.192- 10.16.13.207
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "us1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "MongoIntegrationDebuggerReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "13.208/28", // check for availability of subnet CIDR, we specify 28 to reserve only 16 IPs 10.16.13.208 - 10.16.13.223
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-05e475e8a71d738b5")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true

	defaults.EnableSensuV3ClientEcsService = true
	defaults.EnableSplunkPersistentState = true
	defaults.EnableMongoLogger = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("us1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "mongo-integration-debugger-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}

	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance013180 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013180.EnableEc2instance = true
	MongoReplicaSetInstance013180.Ec2Instance.ImageId = cloudformation.String("ami-05e475e8a71d738b5")
	MongoReplicaSetInstance013180.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance013180.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance013180.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.13.180")
	MongoReplicaSetInstance013180.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013180.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013180.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013180.StopServices = false
	MongoReplicaSetInstance013180.EnableXvdpGp3 = true
	MongoReplicaSetInstance013180.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013180.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013180.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013196 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013196.EnableEc2instance = true
	MongoReplicaSetInstance013196.Ec2Instance.ImageId = cloudformation.String("ami-05e475e8a71d738b5")
	MongoReplicaSetInstance013196.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance013196.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance013196.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.13.196") //primary
	MongoReplicaSetInstance013196.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013196.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013196.StopServices = false
	MongoReplicaSetInstance013196.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013196.EnableXvdpGp3 = true
	MongoReplicaSetInstance013196.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013196.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013196.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance013212 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance013212.EnableEc2instance = true
	MongoReplicaSetInstance013212.Ec2Instance.ImageId = cloudformation.String("ami-05e475e8a71d738b5")
	MongoReplicaSetInstance013212.Ec2Instance.InstanceType = cloudformation.String("r5.large")
	MongoReplicaSetInstance013212.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance013212.Ec2Instance.PrivateIpAddress = cloudformation.String("10.16.13.212")
	MongoReplicaSetInstance013212.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance013212.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance013212.EnableMongoRegistryCache = true
	MongoReplicaSetInstance013212.StopServices = false
	MongoReplicaSetInstance013212.EnableXvdpGp3 = true
	MongoReplicaSetInstance013212.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance013212.MongoContainerTag = "github-sne-6117-16"
	MongoReplicaSetInstance013212.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/us1/Mongo-Integration-Debugger-1", "us1-Mongo-Integration-Debugger-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/us1/Mongo-Integration-Debugger-1", "us1-Mongo-Integration-Debugger-1-Service.json")
}
