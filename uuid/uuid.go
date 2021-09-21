package uuid

import "github.com/bwmarrin/snowflake"

func GetUuid() int64 {
	var i int64 = 1
	for {
		uuid,err := GetSnowFlake(i)
		if err == nil{
			return uuid
		}
	}
}

func GetSnowFlake(machineID int64) (int64,error) {
	node, err := snowflake.NewNode(machineID)
	if err != nil {
		return 0,err
	}
	return node.Generate().Int64(),nil
}
