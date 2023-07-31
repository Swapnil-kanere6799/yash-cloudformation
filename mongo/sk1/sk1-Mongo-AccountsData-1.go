package sk1Mongo

import (
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket"
	"git.wizrocket.net/infra/cloudformation/lib/wizrocket/mongo"
	"github.com/awslabs/goformation/v7/cloudformation"
	"github.com/awslabs/goformation/v7/cloudformation/ecs"
)

func GenerateSk1MongoWiredTigerAccountsData3Dot6Template() {
	sTemplate := mongo.NewStackTemplate()
	serviceTemplate := mongo.NewServiceTemplate()

	sTemplate.Resources["MongoEcsCluster"] = &ecs.Cluster{}
	sTemplate.Resources["MongoVolumeXvdpKmsKey"] = mongo.GetDefaultAWSKmsKeyWithTag()
	sTemplate.Resources["MongoEc2InstanceIamRole"] = mongo.GetDefaultIamRole()
	sTemplate.Resources["MongoEc2InstanceIamPolicy"] = mongo.GetDefaultIamPolicy("sk1")
	sTemplate.Resources["MongoEc2InstanceIamInstanceProfile"] = mongo.GetDefaultIamProfile()

	serviceTemplate.Resources["MongoEcsTaskIamRole"] = mongo.GetTaskExecutionIamRole()
	serviceTemplate.Resources["MongoEcsTaskIamPolicy"] = mongo.GetTaskExecutionIamPolicy("sk1")

	subnetA := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "a",
		Ecc2SubnetLogicalId:    "WiredTigerAccountsMongoReplicaSetSubnetA",
		SubnetCidrBlockSuffix:  "7.48/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.33 - 10.14.6.46
	})
	subnetA.AppendToTemplate(sTemplate)

	subnetB := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "b",
		Ecc2SubnetLogicalId:    "WiredTigerAccountsMongoReplicaSetSubnetB",
		SubnetCidrBlockSuffix:  "7.64/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.49 - 10.14.6.62
	})
	subnetB.AppendToTemplate(sTemplate)

	subnetC := mongo.NewSubnet(mongo.Subnet{
		StackPrefix:            "sk1",
		AvailabilityZoneSuffix: "c",
		Ecc2SubnetLogicalId:    "WiredTigerAccountsMongoReplicaSetSubnetC",
		SubnetCidrBlockSuffix:  "7.96/28", // check for availability of subnet CIDR, we specify 27 to reserve only 32 IPs 10.14.6.65 - 10.14.6.78
	})
	subnetC.AppendToTemplate(sTemplate)

	// We start adding Mongo Instances from here
	defaults := GetDefaultMongoConfiguration()
	defaults.XvdpEc2Volume.Size = cloudformation.Int(64)
	defaults.Ec2Instance.ImageId = cloudformation.String("ami-0a91977023aa59ed0")
	defaults.Ec2Instance.DisableApiTermination = cloudformation.Bool(false)
	defaults.EnableCadvisorArtifactoryRepository = true
	defaults.EnableSplunkPersistentState = true
	defaults.Ec2Instance.SecurityGroupIds = []string{
		cloudformation.ImportValue("sk1-SecurityGroup-MongoInstanceEC2SecurityGroupId"),
	}
	defaults.EcsTaskDefinitionCommand = []string{"--dbpath", "/var/lib/mongo", "--replSet", "wa-rs0", "--logpath", "/var/log/mongodb/mongod.log", "--logappend", "--auth", "--oplogSize", "2048", "--keyFile", "/var/lib/mongodb-keyfile"}
	// Flip this to make all the instances disappear
	defaults.EnableEc2instance = true

	MongoReplicaSetInstance007052 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007052.EnableEc2instance = true
	MongoReplicaSetInstance007052.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007052.Ec2Instance.InstanceType = cloudformation.String("c5.xlarge")
	MongoReplicaSetInstance007052.Ec2InstanceSubnet = subnetA
	MongoReplicaSetInstance007052.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.52") //primary
	MongoReplicaSetInstance007052.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007052.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007052.StopServices = false
	MongoReplicaSetInstance007052.EnableXvdpGp3 = true
	MongoReplicaSetInstance007052.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007052.MongoContainerTag = "bamboo-mongo-task-SNE-39921-amazon-linux-version-6"
	MongoReplicaSetInstance007052.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007052.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007068 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007068.EnableEc2instance = true
	MongoReplicaSetInstance007068.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007068.Ec2Instance.InstanceType = cloudformation.String("c5.xlarge")
	MongoReplicaSetInstance007068.Ec2InstanceSubnet = subnetB
	MongoReplicaSetInstance007068.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.68")
	MongoReplicaSetInstance007068.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007068.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007068.StopServices = false
	MongoReplicaSetInstance007068.EnableXvdpGp3 = true
	MongoReplicaSetInstance007068.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007068.MongoContainerTag = "bamboo-mongo-task-SNE-39921-amazon-linux-version-6"
	MongoReplicaSetInstance007068.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007068.AppendToTemplate(sTemplate, serviceTemplate)

	MongoReplicaSetInstance007100 := mongo.NewMongo(defaults)
	MongoReplicaSetInstance007100.EnableEc2instance = true
	MongoReplicaSetInstance007100.Ec2Instance.ImageId = cloudformation.String("ami-083ddd57181340c47")
	MongoReplicaSetInstance007100.Ec2Instance.InstanceType = cloudformation.String("c5.xlarge")
	MongoReplicaSetInstance007100.Ec2InstanceSubnet = subnetC
	MongoReplicaSetInstance007100.Ec2Instance.PrivateIpAddress = cloudformation.String("10.14.7.100")
	MongoReplicaSetInstance007100.XvdpEc2Volume.Size = cloudformation.Int(64)
	MongoReplicaSetInstance007100.EnableMongoArtifactoryRepository = true
	MongoReplicaSetInstance007100.StopServices = false
	MongoReplicaSetInstance007100.EnableXvdpGp3 = true
	MongoReplicaSetInstance007100.XvdpEc2Volume.Iops = cloudformation.Int(3000)
	MongoReplicaSetInstance007100.MongoContainerTag = "bamboo-mongo-task-SNE-39921-amazon-linux-version-6"
	MongoReplicaSetInstance007100.EnableMongoRegistryCache = true
	MongoReplicaSetInstance007100.AppendToTemplate(sTemplate, serviceTemplate)

	wizrocket.WriteTemplate(sTemplate, "/mongo/sk1/Mongo-AccountsData-1", "sk1-Mongo-AccountsData-1.json")
	wizrocket.WriteTemplate(serviceTemplate, "/mongo/sk1/Mongo-AccountsData-1", "sk1-Mongo-AccountsData-1-Service.json")
}
