import os

if __name__ == "__main__":
    ff = open("out.txt", "w")
    with open("migrate-list", 'r') as f:
        lines = f.readlines()
    for line in lines:
        strs =  line.strip('\n').split(" ")
        albId = strs[0]
        instanceId = strs[1]
        tenantId = strs[2]
        imageId = strs[3]
        zone = strs[4]

        create = "lbs alb-instance-create %s --zone=%s --tenant-id %s --image-id %s" % (albId,zone,tenantId,imageId)
        delete = "lbc alb-instance-delete %s --tenant-id %s" % (instanceId,tenantId)

        print >> ff, "%s\n%s" % (create,delete)